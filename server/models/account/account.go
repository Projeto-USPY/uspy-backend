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

type operation struct {
	ref     *firestore.DocumentRef
	method  string
	payload interface{}

	err error
}

// InsertUser takes the user object and their transcripts and performs all the required database insertions
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

		// register user major
		major := models.Major{Course: data.Course, Specialization: data.Specialization}
		objs = append(objs, db.Object{
			Collection: "users/" + newUser.Hash() + "/majors",
			Doc:        major.Hash(),
			Data:       major,
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
		return errUserExists
	}

	return nil
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

	storedUser.ID = userID
	account.Profile(ctx, storedUser)
}

// Signup performs all the server-side signup operations.
//
// It validates database data, gets and parses user records, creates the user object and sends the verification email
func Signup(ctx *gin.Context, DB db.Env, signupForm *controllers.SignupForm) {
	// check if email already exists in the database
	hashedEmail := utils.SHA256(signupForm.Email)
	query := DB.Client.Collection("users").Where("email", "==", hashedEmail).Limit(1)
	snaps, err := query.Documents(ctx).GetAll()

	if err != nil || len(snaps) != 0 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, views.ErrInvalidEmail)
		return
	}

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

	// create user object
	newUser, userErr := models.NewUser(
		data.Nusp,
		data.Name,
		signupForm.Email,
		signupForm.Password,
		pdf.CreationDate,
	)

	if userErr != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error generating user: %s", userErr.Error()))
		return
	}

	// insert user object into database
	if err := InsertUser(DB, newUser, &data); err != nil {
		if err == errUserExists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, views.ErrInvalidUser)
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error inserting user %s: %s", data.Nusp, err.Error()))
		return
	}

	// send email verification
	if err := sendEmailVerification(signupForm.Email, newUser.Hash()); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to send email verification to user %s; %s", signupForm.Email, err.Error()))
		return
	}

	account.Signup(ctx, newUser.ID, data)
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
	snap, err := DB.Restore("users", utils.SHA256(login.ID))

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
	account.Login(ctx, login.ID, name)
}

// Logout is a dummy method that simply calls the view method that will unset the access token cookie
func Logout(ctx *gin.Context) {
	account.Logout(ctx)
}

// ChangePassword changes the user's password in the database
// This method requires the user to be logged in
func ChangePassword(ctx *gin.Context, DB db.Env, userID string, resetForm *controllers.PasswordChange) {
	snap, err := DB.Restore("users", utils.SHA256(userID))

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
func ResetPassword(ctx *gin.Context, DB db.Env, recovery *controllers.PasswordRecovery) {
	// parse token
	token, _ := utils.ValidateJWT(recovery.Token, config.Env.JWTSecret) // ignoring error because it was already validated in controller
	claims := token.Claims.(jwt.MapClaims)
	userHash := claims["user"].(string)

	var storedUser models.User

	// assert user exists
	if snap, err := DB.Restore("users", userHash); err != nil {
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
func VerifyAccount(ctx *gin.Context, DB db.Env, verification *controllers.AccountVerification) {
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
func Delete(ctx *gin.Context, DB db.Env, userID string) {
	if deleteErr := DB.Client.RunTransaction(DB.Ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		objects := getUserObjects(ctx, DB, tx, userID)

		log.Printf("user is removing their account, total objects affected: %v\n", len(objects))

		for _, obj := range objects {
			if obj.err != nil {
				return obj.err
			}

			var operationErr error
			switch obj.method {
			case "delete":
				operationErr = tx.Delete(obj.ref)
			case "update":
				operationErr = tx.Update(obj.ref, obj.payload.([]firestore.Update))
			}

			if operationErr != nil {
				return fmt.Errorf("could not apply operation on %#v: %s", obj, operationErr.Error())
			}
		}

		return nil
	}); deleteErr != nil {
		ctx.AbortWithError(http.StatusInternalServerError, deleteErr)
		return
	}

	account.Delete(ctx)
}

// getUserObjects gets all documents necessary (and the corresponding operations that must be applied) for a user to remove their account
func getUserObjects(
	ctx context.Context,
	DB db.Env,
	tx *firestore.Transaction,
	userID string,
) []operation {
	objects := make(chan operation)
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
	results := make([]operation, 0)
	go func() {
		for obj := range objects {
			results = append(results, obj)
		}
	}()

	// add userRef to objects
	objects <- operation{
		ref:    userRef,
		method: "delete",
	}

	wg.Wait()
	return results
}

// getScoreObjects gets a list of grade changes, both from the user's final scores, aswell as the subject grades
func getScoreObjects(
	ctx context.Context,
	DB db.Env,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get deletable user final scores and records
	finalScores := userRef.Collection("final_scores")
	scores, err := tx.DocumentRefs(finalScores).GetAll()
	if err != nil {
		objects <- operation{err: errors.New("failed to get final scores from user: " + err.Error())}
		return
	}

	wg.Add(len(scores))
	for _, scoreRef := range scores {
		go func(scoreRef *firestore.DocumentRef) {
			defer wg.Done()

			// insert final scores in channel
			objects <- operation{
				ref:    scoreRef,
				method: "delete",
			}

			// get records
			records := scoreRef.Collection("records")
			recordsDocs, err := tx.DocumentRefs(records).GetAll()
			if err != nil {
				objects <- operation{err: errors.New("failed to get records from user: " + err.Error())}
				return
			}

			wg.Add(len(recordsDocs))
			for _, recordRef := range recordsDocs {
				go func(recordRef *firestore.DocumentRef) {
					defer wg.Done()

					// insert record in channel
					objects <- operation{
						ref:    recordRef,
						method: "delete",
					}

					// get record value
					recordSnap, err := tx.Get(recordRef)
					if err != nil {
						objects <- operation{err: errors.New("failed to convert record to snap: " + err.Error())}
						return
					}

					gradeValue, err := recordSnap.DataAt("grade")
					if err != nil {
						objects <- operation{err: errors.New("failed to get grade value from snap at field grade: " + err.Error())}
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
						objects <- operation{err: errors.New("failed to get queried grades from subject: " + err.Error())}
						return
					}

					// insert subject grades in channel
					for _, ref := range subjectRecordSnaps {
						objects <- operation{
							ref:    ref.Ref,
							method: "delete",
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
	DB db.Env,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get deletable subject reviews and updates
	subjectReviews := userRef.Collection("subject_reviews")
	userReviews, err := tx.DocumentRefs(subjectReviews).GetAll()

	if err != nil {
		objects <- operation{err: errors.New("failed to get subject reviews from user: " + err.Error())}
		return
	}

	wg.Add(len(userReviews))
	for _, reviewRef := range userReviews {
		go func(reviewRef *firestore.DocumentRef) {
			defer wg.Done()

			// insert user review in channel
			objects <- operation{
				ref:    reviewRef,
				method: "delete",
			}

			reviewSnap, err := tx.Get(reviewRef)
			if err != nil {
				objects <- operation{err: errors.New("failed to get review snap: " + err.Error())}
				return
			}

			// lookup subjects that need to be updated
			categories, err := reviewSnap.DataAt("categories") // get existing review categories
			if err != nil {
				objects <- operation{err: errors.New("failed to get review snap at field categories: " + err.Error())}
				return
			}

			subRef := DB.Client.Doc("subjects/" + reviewSnap.Ref.ID)

			// iterate through categories and look which need to have stats decreased
			for k, v := range categories.(map[string]interface{}) {
				// decrement every category review which was true
				if reflect.ValueOf(v).Kind() == reflect.Bool && v.(bool) {
					path := fmt.Sprintf("stats.%s", k)
					objects <- operation{
						ref:     subRef,
						method:  "update",
						payload: []firestore.Update{{Path: path, Value: firestore.Increment(-1)}},
					}
				}
			}

			// number of total reviews must also be decreased
			objects <- operation{
				ref:     subRef,
				method:  "update",
				payload: []firestore.Update{{Path: "stats.total", Value: firestore.Increment(-1)}},
			}
		}(reviewRef)
	}
}

// getCommentObjects gets a list of comment changes
func getCommentObjects(
	ctx context.Context,
	DB db.Env,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get deletable comments
	comments := userRef.Collection("user_comments")
	commentsRefs, err := tx.DocumentRefs(comments).GetAll()
	if err != nil {
		objects <- operation{err: errors.New("failed to get comments from user: " + err.Error())}
		return
	}

	wg.Add(len(commentsRefs))
	for _, userCommentRef := range commentsRefs {
		go func(userCommentRef *firestore.DocumentRef) {
			defer wg.Done()

			// add user comment to objects channel
			objects <- operation{
				ref:    userCommentRef,
				method: "delete",
			}

			// get user comment
			snap, err := tx.Get(userCommentRef)
			if err != nil {
				objects <- operation{err: errors.New("failed to get user comment: " + err.Error())}
				return
			}

			// bind user comment
			var userComment models.UserComment
			if err := snap.DataTo(&userComment); err != nil {
				objects <- operation{err: errors.New("failed to bind user comment: " + err.Error())}
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
			objects <- operation{
				ref:    commentRef,
				method: "delete",
			}
		}(userCommentRef)
	}
}

// getCommentRatingObjects gets a list of comment rating objects
func getCommentRatingObjects(
	ctx context.Context,
	DB db.Env,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get deletable comment ratings
	commentRatings := userRef.Collection("comment_ratings")
	commentRatingsRefs, err := tx.DocumentRefs(commentRatings).GetAll()
	if err != nil {
		objects <- operation{err: errors.New("failed to get comment ratings from user: " + err.Error())}
		return
	}

	wg.Add(len(commentRatingsRefs))
	for _, commentRatingRef := range commentRatingsRefs {
		go func(commentRatingRef *firestore.DocumentRef) {
			defer wg.Done()

			// add comment rating to objects channel
			objects <- operation{
				ref:    commentRatingRef,
				method: "delete",
			}

			// bind comment rating
			var commentRating models.CommentRating
			snap, err := tx.Get(commentRatingRef)
			if err != nil {
				objects <- operation{err: errors.New("failed to get comment rating: " + err.Error())}
				return
			}

			if err := snap.DataTo(&commentRating); err != nil {
				objects <- operation{err: errors.New("failed to bind comment rating: " + err.Error())}
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
				objects <- operation{err: errors.New("failed to get comment from comment rating: " + err.Error())}
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

				objects <- operation{
					ref:     commentSnaps[0].Ref,
					method:  "update",
					payload: []firestore.Update{{Path: path, Value: firestore.Increment(-1)}},
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

				objects <- operation{
					ref:     targetUserCommentRef,
					method:  "update",
					payload: []firestore.Update{{Path: "comment." + path, Value: firestore.Increment(-1)}},
				}
			}

		}(commentRatingRef)
	}
}

// getCommentReportObjects gets a list of report objects
func getCommentReportObjects(
	ctx context.Context,
	DB db.Env,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get deletable comment reports
	commentReports := userRef.Collection("comment_reports")
	commentReportsRefs, err := tx.DocumentRefs(commentReports).GetAll()
	if err != nil {
		objects <- operation{err: errors.New("failed to get comment reports from user: " + err.Error())}
		return
	}

	wg.Add(len(commentReportsRefs))
	for _, commentReportRef := range commentReportsRefs {
		go func(commentReportRef *firestore.DocumentRef) {
			defer wg.Done()

			// add comment report to objects channel
			objects <- operation{
				ref:    commentReportRef,
				method: "delete",
			}

			// bind comment report
			var commentReport models.CommentReport
			snap, err := tx.Get(commentReportRef)
			if err != nil {
				objects <- operation{err: errors.New("failed to get comment report: " + err.Error())}
				return
			}

			if err := snap.DataTo(&commentReport); err != nil {
				objects <- operation{err: errors.New("failed to bind comment report: " + err.Error())}
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
				objects <- operation{err: errors.New("failed to get comment from comment report: " + err.Error())}
				return
			} else if len(commentSnaps) == 0 {
				return
			} else {
				objects <- operation{
					ref:     commentSnaps[0].Ref,
					method:  "update",
					payload: []firestore.Update{{Path: "reports", Value: firestore.Increment(-1)}},
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

				objects <- operation{
					ref:     targetUserCommentRef,
					method:  "update",
					payload: []firestore.Update{{Path: "comment.reports", Value: firestore.Increment(-1)}},
				}
			}

		}(commentReportRef)
	}
}

// getMajorObjects gets a list of major objects
func getMajorObjects(
	ctx context.Context,
	DB db.Env,
	tx *firestore.Transaction,
	wg *sync.WaitGroup,
	objects chan<- operation,
	userRef *firestore.DocumentRef,
) {
	defer wg.Done()

	// get user majors
	majorCol := userRef.Collection("majors")
	majorRefs, err := tx.DocumentRefs(majorCol).GetAll()

	if err != nil {
		objects <- operation{err: errors.New("failed to get majors from user: " + err.Error())}
		return
	}

	for _, major := range majorRefs {
		objects <- operation{
			ref:    major,
			method: "delete",
		}
	}
}
