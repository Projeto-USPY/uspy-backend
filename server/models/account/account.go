// package models
package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/iddigital"
	"github.com/Projeto-USPY/uspy-backend/server/middleware"
	"github.com/Projeto-USPY/uspy-backend/server/views/account"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserExists = errors.New("user is already registered")
)

func InsertUser(DB db.Env, newUser *models.User, data *iddigital.Transcript) error {
	_, err := DB.Restore("users", newUser.Hash())
	if status.Code(err) == codes.NotFound {
		// user is new
		objs := []db.Object{
			{
				Collection: "users",
				Doc:        newUser.Hash(),
				Data:       newUser,
			},
		}

		for _, g := range data.Grades {
			rec := models.Record{
				Grade:     g.Grade,
				Status:    g.Status,
				Frequency: g.Frequency,
				Year:      g.Year,
				Semester:  g.Semester,
			}

			subHash := models.Subject{Code: g.Subject, CourseCode: g.Course, Specialization: g.Specialization}.Hash()

			// store all user records
			objs = append(objs, db.Object{
				Collection: "users/" + newUser.Hash() + "/final_scores/" + subHash + "/records",
				Doc:        rec.Hash(),
				Data:       rec,
			})

			// add grade to "global" grades collection
			gradeObj := models.Record{
				Grade: g.Grade,
			}

			objs = append(objs, db.Object{
				Collection: "subjects/" + subHash + "/grades",
				Data:       gradeObj,
			})
		}

		// write atomically
		if writeErr := DB.BatchWrite(objs); writeErr != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		return ErrUserExists
	}

	return nil
}

// Profile retrieves the user profile from the database
func Profile(ctx *gin.Context, DB db.Env, userID string) {
	var storedUser models.User

	snap, err := DB.Restore("users", utils.SHA256(userID))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user with id %s: %s", userID, err.Error()))
		return
	}
	err = snap.DataTo(&storedUser)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind user %s data to model: %s", userID, err.Error()))
		return
	}

	account.Profile(ctx, storedUser)
}

// Signup inserts a new user into the DB
func Signup(ctx *gin.Context, DB db.Env, signupForm *controllers.SignupForm) {
	// get user records
	cookies := ctx.Request.Cookies()
	resp, err := iddigital.PostAuthCode(signupForm.AccessKey, signupForm.Captcha, cookies)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting pdf from iddigital: %s", err.Error()))
		return
	} else if resp.Header.Get("Content-Type") != "application/pdf" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if pdf := iddigital.NewPDF(resp); pdf.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error converting pdf to text: %s", pdf.Error.Error()))
		return
	} else {
		data, err := pdf.Parse(DB)

		var maxPDFAge float64
		if config.Env.Mode == "dev" {
			maxPDFAge = 24 * 30 // a month
		} else {
			maxPDFAge = 1.0 // an hour
		}

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error parsing pdf: %s", err.Error()))
			return
		} else if time.Since(pdf.CreationDate).Hours() > maxPDFAge {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		newUser, userErr := models.NewUser(
			data.Nusp,
			data.Name,
			signupForm.Password,
			pdf.CreationDate,
		)

		if userErr != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error generating user: %s", userErr.Error()))
			return
		}

		if err := InsertUser(DB, newUser, &data); err != nil {
			if err == ErrUserExists {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}

			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error inserting user %s: %s", data.Nusp, err.Error()))
			return
		}

		account.Signup(ctx, newUser.ID, data)
	}

}

// SignupCaptcha gets the iddigital validation captcha
func SignupCaptcha(ctx *gin.Context) {
	resp, err := iddigital.GetCaptcha()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting captcha from iddigital: %s", err.Error()))
		return
	} else if resp.StatusCode != http.StatusOK {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	account.SignupCaptcha(ctx, resp)
}

// Login performs the user login by comparing the passwordHash and the stored hash
func Login(ctx *gin.Context, DB db.Env, login *controllers.Login) {
	if snap, err := DB.Restore("users", utils.SHA256(login.ID)); err != nil { // get user from database
		if status.Code(err) == codes.NotFound { // if user was not found
			ctx.AbortWithError(http.StatusForbidden, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	} else {
		var storedUser models.User
		if err := snap.DataTo(&storedUser); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err) // some error happened
			return
		}

		// check if password is correct
		if !utils.BcryptCompare(login.Password, storedUser.PasswordHash) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// generate access_token
		if jwtToken, err := middleware.GenerateJWT(login.ID); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error generating jwt for user %s: %s", storedUser.ID, err.Error()))
			return
		} else {
			domain := ctx.MustGet("front_domain").(string)

			// expiration date = 1 month
			secureCookie := !config.Env.IsLocal()
			cookieAge := 0

			// remember this login?
			if login.Remember {
				cookieAge = 30 * 24 * 3600 // 30 days in seconds
			}

			if name, err := utils.AESDecrypt(storedUser.NameHash, config.Env.AESKey); err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error decrypting name hash: %s", err.Error()))
				return
			} else {
				ctx.SetCookie("access_token", jwtToken, cookieAge, "/", domain, secureCookie, true)
				account.Login(ctx, login.ID, name)
			}
		}
	}

}

func Logout(ctx *gin.Context) {
	account.Logout(ctx)
}

// ChangePassword changes the user's password in the database
// This method requires the user to be logged in
func ChangePassword(ctx *gin.Context, DB db.Env, userID string, resetForm *controllers.PasswordChange) {
	if snap, err := DB.Restore("users", utils.SHA256(userID)); err != nil {
		if status.Code(err) == codes.NotFound { // if user was not found
			ctx.AbortWithError(http.StatusForbidden, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	} else {
		var storedUser models.User
		if err := snap.DataTo(&storedUser); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// check if old password is correct
		if !utils.BcryptCompare(resetForm.OldPassword, storedUser.PasswordHash) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// generate new hash
		newHash, err := utils.Bcrypt(resetForm.NewPassword)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error hashing password: %s", err.Error()))
			return
		}

		// update password hash and set userID to be able to find user document
		storedUser.PasswordHash = newHash
		storedUser.ID = userID
		if err := DB.Update(storedUser, "users"); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update password: %s", err.Error()))
			return
		}
	}

	account.ChangePassword(ctx)
}

// ResetPassword resets the user's password in the database
// It differs from ChangePassword because it does not requires an access token
func ResetPassword(ctx *gin.Context, DB db.Env, signupForm *controllers.SignupForm) {
	// get user records
	cookies := ctx.Request.Cookies()
	resp, err := iddigital.PostAuthCode(signupForm.AccessKey, signupForm.Captcha, cookies)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting pdf from iddigital: %s", err.Error()))
		return
	} else if resp.Header.Get("Content-Type") != "application/pdf" {
		// wrong captcha or auth
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if pdf := iddigital.NewPDF(resp); pdf.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error converting pdf to text: %s", pdf.Error.Error()))
		return
	} else {
		data, err := pdf.Parse(DB)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error parsing pdf: %s", err.Error()))
			return
		}

		user := models.User{ID: data.Nusp}
		newHash, err := utils.Bcrypt(signupForm.Password)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error hashing password: %s", err.Error()))
			return
		}

		user.PasswordHash = newHash
		if err := DB.Update(user, "users"); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update password: %s", err.Error()))
			return
		}
	}

	account.ResetPassword(ctx)
}

// Delete deletes the given user (removing all of its traces (grades, reviews, etc)
func Delete(ctx *gin.Context, DB db.Env, userID string) {
	userHash := utils.SHA256(userID)
	userRef := DB.Client.Doc("users/" + userHash)
	records := userRef.Collection("final_scores")
	subjectReviews := userRef.Collection("subject_reviews")

	deleteErr := DB.Client.RunTransaction(DB.Ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Do all the reading first
		recordRefs, subjectReviewsRefs := tx.DocumentRefs(records), tx.DocumentRefs(subjectReviews)

		recordsDocs, err := recordRefs.GetAll()
		if err != nil {
			return fmt.Errorf("could not get final scores: %v", err.Error())
		}

		subjectReviewsDocs, err := subjectReviewsRefs.GetAll()
		if err != nil {
			return fmt.Errorf("could not get subject reviews: %v", err.Error())
		}

		var wg sync.WaitGroup
		channelErr := make(chan error, len(recordsDocs)*100)
		mustDelete := make(chan *firestore.DocumentRef)

		for _, subRef := range recordsDocs {
			wg.Add(1)
			go func(subRef *firestore.DocumentRef) {
				defer wg.Done()

				recordsRef := subRef.Collection("records")
				recordsDocs, err := tx.DocumentRefs(recordsRef).GetAll()
				if err != nil {
					channelErr <- err
				}

				gradesCol := DB.Client.Collection("subjects/" + subRef.ID + "/grades")

				// get grades to remove
				for _, recordRef := range recordsDocs {
					wg.Add(1)
					go func(recordRef *firestore.DocumentRef) {
						defer wg.Done()

						// get value of record
						snap, err := tx.Get(recordRef)
						if err != nil {
							channelErr <- err
						}

						// read final score
						var score models.Record
						if err = snap.DataTo(&score); err != nil {
							channelErr <- err
						}

						// finds one grade in subject/subject_id/grades where the grade is the same as score.Grade
						query := gradesCol.Where("grade", "==", score.Grade).Limit(1)
						gradeDocsToRemove, err := tx.Documents(query).GetAll()

						if err != nil {
							channelErr <- err
						}

						// store grade documents that must be deleted
						for _, gradeSnap := range gradeDocsToRemove {
							mustDelete <- gradeSnap.Ref
						}

						mustDelete <- recordRef
					}(recordRef)
				}
			}(subRef)
		}

		// receive deletables
		deletables := make([]*firestore.DocumentRef, 0, 500)
		doneDelete := make(chan struct{})
		go func() {
			for {
				select {
				case ref := <-mustDelete:
					deletables = append(deletables, ref)
				case <-doneDelete:
					return
				}
			}
		}()

		wg.Wait()
		close(channelErr)
		close(doneDelete)

		for e := range channelErr {
			if e != nil {
				return fmt.Errorf("could not delete grades and records: %v", err.Error())
			}
		}

		// For all review documents
		for _, reviewRef := range subjectReviewsDocs {
			rev, err := tx.Get(reviewRef) // get review snapshot
			if err != nil {
				return fmt.Errorf("could not get review snapshot: %v", err.Error())
			}

			categories, err := rev.DataAt("categories") // get existing review categories
			if err != nil {
				return fmt.Errorf("could not get review categories: %v", err.Error())
			}

			subRef := DB.Client.Doc("subjects/" + reviewRef.ID)

			// For all of the categories
			for k, v := range categories.(map[string]interface{}) {
				// decrement every category review which was true
				if reflect.ValueOf(v).Kind() == reflect.Bool && v.(bool) {
					path := fmt.Sprintf("stats.%s", k)
					err = tx.Update(subRef, []firestore.Update{{Path: path, Value: firestore.Increment(-1)}})
					if err != nil {
						return fmt.Errorf("could not update subject categories: %v", err.Error())
					}
				}
			}
			// decrement number of total reviews
			err = tx.Update(subRef, []firestore.Update{{Path: "stats.total", Value: firestore.Increment(-1)}})
			if err != nil {
				return fmt.Errorf("could not decrement number of total reviews: %v", err.Error())
			}
		}

		// delete stuff
		log.Printf("user %s deleted their account, impacted documents: %d\n", userID, len(deletables))
		for _, d := range deletables {
			if err := tx.Delete(d); err != nil {
				return fmt.Errorf("could not delete grades and records: %v", err.Error())
			}
		}

		return tx.Delete(userRef) // deletes the user
	})

	if deleteErr != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to delete user %s: %s", userID, deleteErr.Error()))
		return
	}

	account.Delete(ctx)
}
