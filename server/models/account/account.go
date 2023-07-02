package account

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/iddigital"
	"github.com/Projeto-USPY/uspy-backend/server/views/account"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUserExists = errors.New("user is already registered")
)

// InsertUser takes the user object and their transcripts and performs all the required database insertions
func InsertUser(DB db.Database, newUser *models.User, data *iddigital.Transcript) error {
	_, err := DB.Restore("users/" + newUser.Hash())
	if status.Code(err) == codes.NotFound {
		// user is new
		objs := []db.BatchObject{
			{
				Collection: "users",
				Doc:        newUser.Hash(),
				WriteData:  newUser,
			},
		}

		// register user major
		major := models.Major{Code: data.Course, Specialization: data.Specialization}
		objs = append(objs, db.BatchObject{
			Collection: "users/" + newUser.Hash() + "/majors",
			Doc:        major.Hash(),
			WriteData:  major,
		})

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
			objs = append(objs, db.BatchObject{
				Collection: "users/" + newUser.Hash() + "/final_scores/" + subHash + "/records",
				Doc:        rec.Hash(),
				WriteData:  rec,
			})

			// add grade to "global" grades collection
			gradeObj := models.Record{
				Grade: g.Grade,
			}

			objs = append(objs, db.BatchObject{
				Collection: "subjects/" + subHash + "/grades",
				WriteData:  gradeObj,
			})
		}

		// write atomically
		if writeErr := DB.BatchWrite(objs); writeErr != nil {
			return writeErr
		}
	} else if err != nil {
		return err
	} else {
		return errUserExists
	}

	return nil
}

// UpdateUser takes a new transcript and updates the user stored data
//
// It does a diff operation on the already stored grade transcript, adding new final scores to the user's data and updating subject records
func UpdateUser(ctx context.Context, DB db.Database, data *iddigital.Transcript, userID string, updateTime time.Time) error {
	userHash := utils.SHA256(userID)

	return DB.Client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		ops := make([]db.Operation, 0)
		userRef := DB.Client.Doc(fmt.Sprintf(
			"users/%s",
			userHash,
		))

		// lookup user to get trascript years map and do diff
		snap, err := tx.Get(userRef)
		if err != nil {
			return err
		}

		var storedUser models.User
		if err := snap.DataTo(&storedUser); err != nil {
			return err
		}

		if storedUser.TranscriptYears != nil { // merge maps
			for year, semesters := range data.TranscriptYears {
				if _, ok := storedUser.TranscriptYears[year]; !ok { // add whole new year to transcript
					storedUser.TranscriptYears[year] = semesters
				} else { // merge arrays
					storedUser.TranscriptYears[year] = append(storedUser.TranscriptYears[year], semesters...) // merge
					storedUser.TranscriptYears[year] = utils.UniqueInts(storedUser.TranscriptYears[year])     // unique
				}
			}
		} else { // maybe user is outdated. This was added after a change in signup
			storedUser.TranscriptYears = data.TranscriptYears
		}

		// update user transcript years map and user last update time
		ops = append(ops, db.Operation{
			Ref:    userRef,
			Method: "update",
			Payload: []firestore.Update{
				{
					Path:  "last_update",
					Value: updateTime,
				},
				{
					Path:  "transcript_years",
					Value: storedUser.TranscriptYears,
				},
			},
		})

		major := models.Major{
			Code:           data.Course,
			Specialization: data.Specialization,
		}

		majorRef := DB.Client.Doc(fmt.Sprintf(
			"users/%s/majors/%s",
			userHash,
			major.Hash(),
		))

		// set user major
		ops = append(ops, db.Operation{
			Ref:     majorRef,
			Method:  "set",
			Payload: major,
		})

		newRecords := 0

		// lookup user records that are not yet stored
		recordRefs := make([]*firestore.DocumentRef, 0, len(data.Grades))
		for _, grade := range data.Grades {
			subject := models.Subject{Code: grade.Subject, CourseCode: grade.Course, Specialization: grade.Specialization}

			ref := DB.Client.Doc(fmt.Sprintf(
				"users/%s/final_scores/%s/records/%s",
				userHash,
				subject.Hash(),
				grade.Hash(),
			))

			recordRefs = append(recordRefs, ref)
		}

		snaps, err := tx.GetAll(recordRefs)
		if err != nil {
			return err
		}

		for i := 0; i < len(snaps); i++ { // assuming Get All returns snaps in the same order as refs passed
			snap := snaps[i]
			grade := data.Grades[i]

			if !snap.Exists() {
				newRecords++
				// record must be added
				// add operation to array
				ops = append(ops, db.Operation{
					Ref:     snap.Ref,
					Payload: grade,
					Method:  "set",
				})

				// subject grade must be added
				subject := models.Subject{Code: grade.Subject, CourseCode: grade.Course, Specialization: grade.Specialization}

				// create new ref
				gradeRef := DB.Client.Collection("subjects/" + subject.Hash() + "/grades").NewDoc()
				subjectGradeDoc := models.Record{
					Grade: grade.Grade,
				}

				// add operation to array
				ops = append(ops, db.Operation{
					Ref:     gradeRef,
					Payload: subjectGradeDoc,
					Method:  "set",
				})
			}
		}

		log.WithFields(log.Fields{
			"new_records": newRecords,
			"num_ops":     len(ops),
		}).Debug("applying update operation")

		// apply operations
		return db.ApplyConcurrentOperationsInTransaction(tx, ops)
	})
}

func sendPasswordRecoveryEmail(email, userHash string) error {
	if config.Env.IsLocal() || config.Env.IsDev() {
		return nil
	}

	token, err := utils.GenerateJWT(map[string]interface{}{
		"type":      "password_reset",
		"user":      userHash,
		"timestamp": time.Now(),
	}, config.Env.JWTSecret)

	if err != nil {
		return err
	}

	var host string
	if config.Env.IsDev() {
		host = "frontdev.uspy.me"
	} else {
		host = "uspy.me"
	}

	url := fmt.Sprintf(`https://%s/account/password_reset?token=%s`, host, token)
	content := fmt.Sprintf(config.PasswordRecoveryContent, url)
	return config.Env.Send(email, config.PasswordRecoverySubject, content)
}

func sendEmailVerification(email, userHash string) error {
	if config.Env.IsLocal() || config.Env.IsDev() {
		return nil
	}

	emailHash := utils.SHA256(email)
	token, err := utils.GenerateJWT(map[string]interface{}{
		"type":      "email_verification",
		"user":      userHash,
		"email":     emailHash,
		"timestamp": time.Now(),
	}, config.Env.JWTSecret)

	if err != nil {
		return err
	}

	var host string
	if config.Env.IsDev() {
		host = "frontdev.uspy.me"
	} else {
		host = "uspy.me"
	}

	url := fmt.Sprintf(`https://%s/account/verify?token=%s`, host, token)
	content := fmt.Sprintf(config.VerificationContent, url)
	return config.Env.Send(email, config.VerificationSubject, content)
}

// PreSignup performs the server-side signup operations that can be done before the user is created.
//
// It validates the access key, fetches the PDF, parses it and inserts the data into the database
func PreSignup(ctx *gin.Context, DB db.Database, authForm *controllers.AuthForm) {
	// fetch pdf
	resp, err := iddigital.GetPDF(authForm.AccessKey)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting pdf from iddigital: %s", err.Error()))
		return
	} else if resp.StatusCode == http.StatusBadRequest {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, views.ErrInvalidAuthCode)
		return
	} else if resp.StatusCode != http.StatusOK {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, views.ErrOther)
		return
	}

	// parse pdf
	pdf := iddigital.NewPDF(resp)
	transcript, err := pdf.Parse(DB)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error parsing pdf: %s", err.Error()))
		return
	}

	var maxPDFAge float64
	if config.Env.IsDev() || config.Env.IsLocal() {
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
		transcript.Nusp,
		transcript.Name,
		pdf.CreationDate,
		transcript.TranscriptYears,
	)

	if userErr != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error generating user: %s", userErr.Error()))
		return
	}

	// set signup token in response
	// this is used to complete the signup process
	token, err := utils.GenerateJWT(map[string]interface{}{
		"type":      "signup",
		"timestamp": time.Now(),
		"user":      newUser.Hash(),
	}, config.Env.JWTSecret)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error generating signup token: %s", err.Error()))
		return
	}

	account.PreSignup(ctx, token)

	// insert data into database in batch write
	//
	// data to insert:
	// - user (pending verification)
	// - grades
	// - final scores

	// Note that in dev/local environment, verification is skipped and thus this should not be required when completing signup
	if err := InsertUser(DB, newUser, &transcript); err != nil {
		if err == errUserExists {
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error inserting user %s: %s", transcript.Nusp, err.Error()))
		return
	}
}

// CompleteSignup performs the server-side signup operations that can only be done after the user is created.
//
// It validates the signup token, checks if the email is already in use and sends the verification email
func CompleteSignup(ctx *gin.Context, DB db.Database, signupForm *controllers.CompleteSignupForm) {
	// check if email already exists in the database
	hashedEmail := utils.SHA256(signupForm.Email)
	query := DB.Client.Collection("users").Where("email", "==", hashedEmail).Limit(1)
	snaps, err := query.Documents(ctx).GetAll()

	if err != nil || len(snaps) != 0 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, views.ErrInvalidEmail)
		return
	}

	// unwrap signup token
	token, err := utils.ValidateJWT(signupForm.SignupToken, config.Env.JWTSecret)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("error validating signup token: %s", err.Error()))
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userHash := claims["user"].(string)

	// check if user exists
	// if user exists, check if the user is pending verification
	// if user is not pending verification, abort with error

	// get user from database
	snap, err := DB.Restore("users/" + userHash)
	if err != nil {
		if status.Code(err) == codes.NotFound { // if user was not found
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, views.ErrInvalidUser)
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting user: %s", err.Error()))
		}
	}

	var user models.User
	if err := snap.DataTo(&user); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting user: %s", err.Error()))
		return
	}

	// if in production, check if user is pending verification
	//
	// locally or in dev, skip this check
	if config.Env.IsProd() {
		if user.Verified {
			ctx.AbortWithStatusJSON(http.StatusForbidden, views.ErrInvalidUser)
			return
		}
	}

	// complete signup
	err = models.CompleteSignup(DB, userHash, "users", signupForm.Email, signupForm.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error completing signup: %s", err.Error()))
		return
	}

	// send verification email
	if err := sendEmailVerification(signupForm.Email, userHash); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error sending verification email: %s", err.Error()))
		return
	}

	account.CompleteSignup(ctx)
}

// Login performs the user login by comparing the passwordHash and the stored hash
func Login(ctx *gin.Context, DB db.Database, login *controllers.Login) {
	snap, err := DB.Restore("users/" + utils.SHA256(login.ID))

	if err != nil { // get user from database
		if status.Code(err) == codes.NotFound { // if user was not found
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var storedUser models.User
	if err := snap.DataTo(&storedUser); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err) // some error happened
		return
	}

	// check if password is correct
	if !utils.BcryptCompare(login.Password, storedUser.PasswordHash) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, views.ErrInvalidCredentials)
		return
	}

	// check if user has verified their email
	if !storedUser.Verified {
		ctx.AbortWithStatusJSON(http.StatusForbidden, views.ErrUnverifiedUser)
		return
	}

	// check if user is banned
	if storedUser.Banned {
		ctx.AbortWithStatusJSON(http.StatusForbidden, views.ErrBannedUser)
		return
	}

	// generate access_token
	jwtToken, err := utils.GenerateJWT(map[string]interface{}{
		"type":      "access",
		"user":      login.ID,
		"timestamp": time.Now().Unix(),
	}, config.Env.JWTSecret)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error generating jwt for user %s: %s", storedUser.ID, err.Error()))
		return
	}

	domain := ctx.MustGet("front_domain").(string)

	// expiration date = 1 month
	secureCookie := !config.Env.IsLocal()
	cookieAge := 0

	// remember this login?
	if login.Remember {
		cookieAge = 30 * 24 * 3600 // 30 days in seconds
	}

	name, err := utils.AESDecrypt(storedUser.NameHash, config.Env.AESKey)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error decrypting name hash: %s", err.Error()))
		return
	}

	ctx.SetCookie("access_token", jwtToken, cookieAge, "/", domain, secureCookie, true)
	account.Login(ctx, login.ID, name, storedUser.LastUpdate)
}

// Logout is a dummy method that simply calls the view method that will unset the access token cookie
func Logout(ctx *gin.Context) {
	account.Logout(ctx)
}

// ChangePassword changes the user's password in the database
// This method requires the user to be logged in
func ChangePassword(ctx *gin.Context, DB db.Database, userID string, resetForm *controllers.PasswordChange) {
	snap, err := DB.Restore("users/" + utils.SHA256(userID))

	if err != nil {
		if status.Code(err) == codes.NotFound { // if user was not found
			ctx.AbortWithError(http.StatusForbidden, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var storedUser models.User
	if err := snap.DataTo(&storedUser); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// check if old password is correct
	if !utils.BcryptCompare(resetForm.OldPassword, storedUser.PasswordHash) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, views.ErrWrongPassword)
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

	account.ChangePassword(ctx)
}

// ResetPassword resets the user's password in the database
//
// It differs from ChangePassword because it does not requires an access token
func ResetPassword(ctx *gin.Context, DB db.Database, recovery *controllers.PasswordRecovery) {
	// parse token
	token, _ := utils.ValidateJWT(recovery.Token, config.Env.JWTSecret) // ignoring error because it was already validated in controller
	claims := token.Claims.(jwt.MapClaims)
	userHash := claims["user"].(string)

	var storedUser models.User

	// assert user exists
	if snap, err := DB.Restore("users/" + userHash); err != nil {
		if status.Code(err) == codes.NotFound { // if user was not found
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	} else if err := snap.DataTo(&storedUser); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// generate new hash
	newHash, err := utils.Bcrypt(recovery.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error hashing password: %s", err.Error()))
		return
	}

	// update password hash and set userID to be able to find user document
	storedUser.PasswordHash = newHash
	storedUser.IDHash = userHash
	if err := DB.Update(storedUser, "users"); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update password: %s", err.Error()))
		return
	}

	// update user password
	account.ResetPassword(ctx)
}

// VerifyAccount sets the user's email as verified
func VerifyAccount(ctx *gin.Context, DB db.Database, verification *controllers.AccountVerification) {
	token, _ := utils.ValidateJWT(verification.Token, config.Env.JWTSecret) // ignoring error because it was already validated in controller
	claims := token.Claims.(jwt.MapClaims)

	userHash := claims["user"].(string)

	// verify user
	user := models.User{IDHash: userHash, Verified: true}

	if err := DB.Update(user, "users"); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	account.VerifyAccount(ctx)
}

// Delete deletes the given user (removing all of its traces (grades, reviews, etc))
//
// This function wraps the document lookups and simply performs the necessary database operations
func Delete(ctx *gin.Context, DB db.Database, userID string) {
	if deleteErr := DB.Client.RunTransaction(DB.Ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		objects := getUserObjects(ctx, DB, tx, userID)

		log.WithFields(log.Fields{
			"affected_objects": len(objects),
		}).Debug("user is removing their account")
		return db.ApplyConcurrentOperationsInTransaction(tx, objects)
	}); deleteErr != nil {
		ctx.AbortWithError(http.StatusInternalServerError, deleteErr)
		return
	}

	account.Delete(ctx)
}

// Update updates a user's profile with a new grade transcript
func Update(ctx *gin.Context, DB db.Database, userID string, updateForm *controllers.UpdateForm) {
	// get user records
	resp, err := iddigital.GetPDF(updateForm.AccessKey)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting pdf from iddigital: %s", err.Error()))
		return
	} else if resp.StatusCode == http.StatusBadRequest {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, views.ErrInvalidAuthCode)
		return
	} else if resp.StatusCode != http.StatusOK {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, views.ErrOther)
		return
	}

	// parse transcript
	pdf := iddigital.NewPDF(resp)
	if pdf.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error converting pdf to text: %s", pdf.Error.Error()))
		return
	}

	data, err := pdf.Parse(DB)

	var maxPDFAge float64
	if config.Env.IsDev() || config.Env.IsLocal() {
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

	if userID != data.Nusp {
		ctx.AbortWithStatusJSON(http.StatusForbidden, views.ErrInvalidUpdate)
		return
	}

	if err := UpdateUser(ctx, DB, &data, userID, pdf.CreationDate); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error updating user: %s", err.Error()))
		return
	}

	account.Update(ctx)
}

// getUserObjects gets all documents necessary (and the corresponding operations that must be applied) for a user to remove their account
func getUserObjects(
	ctx context.Context,
	DB db.Database,
	tx *firestore.Transaction,
	userID string,
) []db.Operation {
	objects := make(chan db.Operation)
	defer close(objects)

	userHash := utils.SHA256(userID)
	userRef := DB.Client.Doc("users/" + userHash)

	var wg sync.WaitGroup
	wg.Add(6)

	go getScoreObjects(ctx, DB, tx, &wg, objects, userRef)
	go getReviewObjects(ctx, DB, tx, &wg, objects, userRef)
	go getCommentObjects(ctx, DB, tx, &wg, objects, userRef)
	go getCommentRatingObjects(ctx, DB, tx, &wg, objects, userRef)
	go getCommentReportObjects(ctx, DB, tx, &wg, objects, userRef)
	go getMajorObjects(ctx, DB, tx, &wg, objects, userRef)

	// get collected objects and append to array
	results := make([]db.Operation, 0)
	go func() {
		for obj := range objects {
			results = append(results, obj)
		}
	}()

	// add userRef to objects
	objects <- db.Operation{
		Ref:    userRef,
		Method: "delete",
	}

	wg.Wait()
	return results
}

// getScoreObjects gets a list of grade changes, both from the user's final scores, aswell as the subject grades
func getScoreObjects(
	ctx context.Context,
	DB db.Database,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- db.Operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get deletable user final scores and records
	finalScores := userRef.Collection("final_scores")
	scores, err := tx.DocumentRefs(finalScores).GetAll()
	if err != nil {
		objects <- db.Operation{Err: errors.New("failed to get final scores from user: " + err.Error())}
		return
	}

	wg.Add(len(scores))
	for _, scoreRef := range scores {
		go func(scoreRef *firestore.DocumentRef) {
			defer wg.Done()

			// insert final scores in channel
			objects <- db.Operation{
				Ref:    scoreRef,
				Method: "delete",
			}

			// get records
			records := scoreRef.Collection("records")
			recordsDocs, err := tx.DocumentRefs(records).GetAll()
			if err != nil {
				objects <- db.Operation{Err: errors.New("failed to get records from user: " + err.Error())}
				return
			}

			wg.Add(len(recordsDocs))
			for _, recordRef := range recordsDocs {
				go func(recordRef *firestore.DocumentRef) {
					defer wg.Done()

					// insert record in channel
					objects <- db.Operation{
						Ref:    recordRef,
						Method: "delete",
					}

					// get record value
					recordSnap, err := tx.Get(recordRef)
					if err != nil {
						objects <- db.Operation{Err: errors.New("failed to convert record to snap: " + err.Error())}
						return
					}

					gradeValue, err := recordSnap.DataAt("grade")
					if err != nil {
						objects <- db.Operation{Err: errors.New("failed to get grade value from snap at field grade: " + err.Error())}
						return
					}

					// query object with same value in subject grades
					grades := DB.Client.Collection(fmt.Sprintf(
						"subjects/%s/grades",
						scoreRef.ID,
					))

					query := grades.Where("grade", "==", gradeValue).Limit(1)
					subjectRecordSnaps, err := tx.Documents(query).GetAll()
					if err != nil {
						objects <- db.Operation{Err: errors.New("failed to get queried grades from subject: " + err.Error())}
						return
					}

					// insert subject grades in channel
					for _, ref := range subjectRecordSnaps {
						objects <- db.Operation{
							Ref:    ref.Ref,
							Method: "delete",
						}
					}
				}(recordRef)
			}
		}(scoreRef)
	}
}

// getReviewObjects gets a list of review changes, both from the user's reviews, aswell as the subject grades that must have their stats updated
func getReviewObjects(
	ctx context.Context,
	DB db.Database,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- db.Operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get deletable subject reviews and updates
	subjectReviews := userRef.Collection("subject_reviews")
	userReviews, err := tx.DocumentRefs(subjectReviews).GetAll()

	if err != nil {
		objects <- db.Operation{Err: errors.New("failed to get subject reviews from user: " + err.Error())}
		return
	}

	wg.Add(len(userReviews))
	for _, reviewRef := range userReviews {
		go func(reviewRef *firestore.DocumentRef) {
			defer wg.Done()

			// insert user review in channel
			objects <- db.Operation{
				Ref:    reviewRef,
				Method: "delete",
			}

			reviewSnap, err := tx.Get(reviewRef)
			if err != nil {
				objects <- db.Operation{Err: errors.New("failed to get review snap: " + err.Error())}
				return
			}

			// lookup subjects that need to be updated
			categories, err := reviewSnap.DataAt("categories") // get existing review categories
			if err != nil {
				objects <- db.Operation{Err: errors.New("failed to get review snap at field categories: " + err.Error())}
				return
			}

			subRef := DB.Client.Doc("subjects/" + reviewSnap.Ref.ID)

			// iterate through categories and look which need to have stats decreased
			for k, v := range categories.(map[string]interface{}) {
				// decrement every category review which was true
				if reflect.ValueOf(v).Kind() == reflect.Bool && v.(bool) {
					path := fmt.Sprintf("stats.%s", k)
					objects <- db.Operation{
						Ref:     subRef,
						Method:  "update",
						Payload: []firestore.Update{{Path: path, Value: firestore.Increment(-1)}},
					}
				}
			}

			// number of total reviews must also be decreased
			objects <- db.Operation{
				Ref:     subRef,
				Method:  "update",
				Payload: []firestore.Update{{Path: "stats.total", Value: firestore.Increment(-1)}},
			}
		}(reviewRef)
	}
}

// getCommentObjects gets a list of comment changes
func getCommentObjects(
	ctx context.Context,
	DB db.Database,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- db.Operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get deletable comments
	comments := userRef.Collection("user_comments")
	commentsRefs, err := tx.DocumentRefs(comments).GetAll()
	if err != nil {
		objects <- db.Operation{Err: errors.New("failed to get comments from user: " + err.Error())}
		return
	}

	wg.Add(len(commentsRefs))
	for _, userCommentRef := range commentsRefs {
		go func(userCommentRef *firestore.DocumentRef) {
			defer wg.Done()

			// add user comment to objects channel
			objects <- db.Operation{
				Ref:    userCommentRef,
				Method: "delete",
			}

			// get user comment
			snap, err := tx.Get(userCommentRef)
			if err != nil {
				objects <- db.Operation{Err: errors.New("failed to get user comment: " + err.Error())}
				return
			}

			// bind user comment
			var userComment models.UserComment
			if err := snap.DataTo(&userComment); err != nil {
				objects <- db.Operation{Err: errors.New("failed to bind user comment: " + err.Error())}
				return
			}

			// get comment from offerings subcollection
			subject := models.Subject{Code: userComment.Subject, CourseCode: userComment.Course, Specialization: userComment.Specialization}
			commentRef := DB.Client.Doc(fmt.Sprintf(
				"subjects/%s/offerings/%s/comments/%s",
				subject.Hash(),
				userComment.ProfessorHash,
				userRef.ID,
			))

			// add comment to objects channel
			objects <- db.Operation{
				Ref:    commentRef,
				Method: "delete",
			}
		}(userCommentRef)
	}
}

// getCommentRatingObjects gets a list of comment rating objects
func getCommentRatingObjects(
	ctx context.Context,
	DB db.Database,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- db.Operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get deletable comment ratings
	commentRatings := userRef.Collection("comment_ratings")
	commentRatingsRefs, err := tx.DocumentRefs(commentRatings).GetAll()
	if err != nil {
		objects <- db.Operation{Err: errors.New("failed to get comment ratings from user: " + err.Error())}
		return
	}

	wg.Add(len(commentRatingsRefs))
	for _, commentRatingRef := range commentRatingsRefs {
		go func(commentRatingRef *firestore.DocumentRef) {
			defer wg.Done()

			// add comment rating to objects channel
			objects <- db.Operation{
				Ref:    commentRatingRef,
				Method: "delete",
			}

			// bind comment rating
			var commentRating models.CommentRating
			snap, err := tx.Get(commentRatingRef)
			if err != nil {
				objects <- db.Operation{Err: errors.New("failed to get comment rating: " + err.Error())}
				return
			}

			if err := snap.DataTo(&commentRating); err != nil {
				objects <- db.Operation{Err: errors.New("failed to bind comment rating: " + err.Error())}
				return
			}

			// lookup rated comment to update upvotes/downvotes count
			subject := models.Subject{Code: commentRating.Subject, CourseCode: commentRating.Course, Specialization: commentRating.Specialization}
			commentsCol := DB.Client.Collection(fmt.Sprintf(
				"subjects/%s/offerings/%s/comments",
				subject.Hash(),
				commentRating.ProfessorHash,
			))

			query := commentsCol.Where("id", "==", commentRating.ID)
			if commentSnaps, err := tx.Documents(query).GetAll(); err != nil {
				objects <- db.Operation{Err: errors.New("failed to get comment from comment rating: " + err.Error())}
				return
			} else if len(commentSnaps) == 0 { // this comment does not exist anymore
				return
			} else {
				var path string

				if commentRating.Upvote { // decrease original comment upvote count
					path = "upvotes"
				} else { // decrease original comment downvote count
					path = "downvotes"
				}

				objects <- db.Operation{
					Ref:     commentSnaps[0].Ref,
					Method:  "update",
					Payload: []firestore.Update{{Path: path, Value: firestore.Increment(-1)}},
				}

				targetUserComment := models.UserComment{
					ProfessorHash:  commentRating.ProfessorHash,
					Subject:        commentRating.Subject,
					Course:         commentRating.Course,
					Specialization: commentRating.Specialization,
				}

				// user comment should receive same update
				targetUserCommentRef := DB.Client.Doc(fmt.Sprintf(
					"users/%s/user_comments/%s",
					commentSnaps[0].Ref.ID,
					targetUserComment.Hash(),
				))

				objects <- db.Operation{
					Ref:     targetUserCommentRef,
					Method:  "update",
					Payload: []firestore.Update{{Path: "comment." + path, Value: firestore.Increment(-1)}},
				}
			}

		}(commentRatingRef)
	}
}

// getCommentReportObjects gets a list of report objects
func getCommentReportObjects(
	ctx context.Context,
	DB db.Database,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- db.Operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get deletable comment reports
	commentReports := userRef.Collection("comment_reports")
	commentReportsRefs, err := tx.DocumentRefs(commentReports).GetAll()
	if err != nil {
		objects <- db.Operation{Err: errors.New("failed to get comment reports from user: " + err.Error())}
		return
	}

	wg.Add(len(commentReportsRefs))
	for _, commentReportRef := range commentReportsRefs {
		go func(commentReportRef *firestore.DocumentRef) {
			defer wg.Done()

			// add comment report to objects channel
			objects <- db.Operation{
				Ref:    commentReportRef,
				Method: "delete",
			}

			// bind comment report
			var commentReport models.CommentReport
			snap, err := tx.Get(commentReportRef)
			if err != nil {
				objects <- db.Operation{Err: errors.New("failed to get comment report: " + err.Error())}
				return
			}

			if err := snap.DataTo(&commentReport); err != nil {
				objects <- db.Operation{Err: errors.New("failed to bind comment report: " + err.Error())}
				return
			}

			// lookup reported comment to update reports count
			subject := models.Subject{Code: commentReport.Subject, CourseCode: commentReport.Course, Specialization: commentReport.Specialization}
			commentsCol := DB.Client.Collection(fmt.Sprintf(
				"subjects/%s/offerings/%s/comments",
				subject.Hash(),
				commentReport.ProfessorHash,
			))

			query := commentsCol.Where("id", "==", commentReport.ID)
			if commentSnaps, err := tx.Documents(query).GetAll(); err != nil {
				objects <- db.Operation{Err: errors.New("failed to get comment from comment report: " + err.Error())}
				return
			} else if len(commentSnaps) == 0 {
				return
			} else {
				objects <- db.Operation{
					Ref:     commentSnaps[0].Ref,
					Method:  "update",
					Payload: []firestore.Update{{Path: "reports", Value: firestore.Increment(-1)}},
				}

				targetUserComment := models.UserComment{
					ProfessorHash:  commentReport.ProfessorHash,
					Subject:        commentReport.Subject,
					Course:         commentReport.Course,
					Specialization: commentReport.Specialization,
				}

				// user comment should receive same update
				targetUserCommentRef := DB.Client.Doc(fmt.Sprintf(
					"users/%s/user_comments/%s",
					commentSnaps[0].Ref.ID,
					targetUserComment.Hash(),
				))

				objects <- db.Operation{
					Ref:     targetUserCommentRef,
					Method:  "update",
					Payload: []firestore.Update{{Path: "comment.reports", Value: firestore.Increment(-1)}},
				}
			}

		}(commentReportRef)
	}
}

// getMajorObjects gets a list of major objects
func getMajorObjects(
	ctx context.Context,
	DB db.Database,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- db.Operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get user majors
	majorCol := userRef.Collection("majors")
	majorRefs, err := tx.DocumentRefs(majorCol).GetAll()

	if err != nil {
		objects <- db.Operation{Err: errors.New("failed to get majors from user: " + err.Error())}
		return
	}

	for _, major := range majorRefs {
		objects <- db.Operation{
			Ref:    major,
			Method: "delete",
		}
	}
}
