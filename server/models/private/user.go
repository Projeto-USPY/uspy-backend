// package models
package private

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	db_utils "github.com/Projeto-USPY/uspy-backend/db/utils"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/private"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetSubjectGrade is the model implementation for /server/controller/private/user.GetSubjectGrade
func GetSubjectGrade(ctx *gin.Context, DB db.Env, userID string, sub *controllers.Subject) {
	user, model := models.User{ID: userID}, models.NewSubjectFromController(sub)
	userHash, subHash := user.Hash(), model.Hash()

	col, err := DB.RestoreCollection("users/" + userHash + "/final_scores/" + subHash + "/records")
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find records: %s", err.Error()))
		return
	} else if len(col) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	best := models.Record{}
	for _, s := range col {
		var fs models.Record
		err := s.DataTo(&fs)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not bind final score: %s", err.Error()))
			return
		}

		if fs.Grade > best.Grade {
			best = fs
		} else if fs.Grade == best.Grade && fs.Year > best.Year {
			best = fs
		}
	}

	private.GetSubjectGrade(ctx, &best)
}

// GetSubjectReview is the model implementation for /server/controller/private/user.GetSubjectReview
func GetSubjectReview(ctx *gin.Context, DB db.Env, userID string, sub *controllers.Subject) {
	user, model := models.User{ID: userID}, models.NewSubjectFromController(sub)
	userHash, subHash := user.Hash(), model.Hash()
	review := models.SubjectReview{}

	err := db_utils.CheckSubjectPermission(DB, userHash, subHash)
	if err != nil {
		if err == db_utils.ErrSubjectNotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject %v: %s", model, err.Error()))
			return
		}

		if err == db_utils.ErrNoPermission {
			ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("user %v has no permission to get review: %s", userID, err.Error()))
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting subject review: %s", err.Error()))
		return
	}

	snap, err := DB.Restore("users/"+userHash+"/subject_reviews", subHash)

	if err != nil { // user has not reviewed subject
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject review for %v and user %v: %s", model, userID, err.Error()))
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not get subject review: %s", err.Error()))
		return
	}

	if err := snap.DataTo(&review); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not bind subject review: %s", err.Error()))
		return
	}

	private.GetSubjectReview(ctx, &review)
}

// UpdateSubjectReview is the model implementation for /server/controller/private/user.UpdateSubjectReview
func UpdateSubjectReview(ctx *gin.Context, DB db.Env, userID string, review *controllers.SubjectReview) {
	userHash, model := models.User{ID: userID}.Hash(), models.NewSubjectReviewFromController(review)
	err := db_utils.CheckSubjectPermission(DB, userHash, model.Hash())
	if err != nil {
		if err == db_utils.ErrSubjectNotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject %v: %s", model, err.Error()))
			return
		}

		if err == db_utils.ErrNoPermission {
			ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("user %v has no permission to get review: %s", userID, err.Error()))
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error checking subject permission: %s", err.Error()))
		return
	}

	revRef := DB.Client.Doc("users/" + userHash + "/subject_reviews/" + model.Hash())
	subRef := DB.Client.Doc("subjects/" + model.Hash())

	err = DB.Client.RunTransaction(DB.Ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		rev, _ := tx.Get(revRef) // get existing review

		if rev != nil && rev.Exists() { // user has already reviewed subject so we must remove it and propagate
			if err != nil {
				return err
			}

			categories, err := rev.DataAt("categories") // get existing review categories
			if err != nil {
				return err
			}

			for k, v := range categories.(map[string]interface{}) {
				// decrement every category review which was true
				if reflect.ValueOf(v).Kind() == reflect.Bool && v.(bool) {
					path := fmt.Sprintf("stats.%s", k)
					err = tx.Update(subRef, []firestore.Update{{Path: path, Value: firestore.Increment(-1)}})
					if err != nil {
						return err
					}
				}
			}

			// decrement number of total reviews
			err = tx.Update(subRef, []firestore.Update{{Path: "stats.total", Value: firestore.Increment(-1)}})
			if err != nil {
				return err
			}
		}

		// add new review (overwrites if existing)
		err := tx.Set(revRef, model)
		if err != nil {
			return err
		}

		// update subject stats with new review
		for k, v := range model.Review {
			if reflect.ValueOf(v).Kind() == reflect.Bool && v.(bool) {
				path := fmt.Sprintf("stats.%s", k)
				err = tx.Update(subRef, []firestore.Update{{Path: path, Value: firestore.Increment(1)}})
				if err != nil {
					return err
				}
			}
		}

		// increment total of reviews
		return tx.Update(subRef, []firestore.Update{{Path: "stats.total", Value: firestore.Increment(1)}})
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error updating subject review: %s", err.Error()))
		return
	}

	private.UpdateSubjectReview(ctx)
}
