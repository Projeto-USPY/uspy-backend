package private

import (
	"cloud.google.com/go/firestore"
	"errors"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"golang.org/x/net/context"
	"reflect"
)

func checkReviewPermission(DB db.Env, userHash, subHash string) error {
	col, err := DB.RestoreCollection("users/" + userHash + "/final_scores/" + subHash + "/records")

	if len(col) == 0 || err != nil { // user has not done subject
		return errors.New("user has not done subject")
	}

	return nil
}

func GetSubjectReview(DB db.Env, user entity.User, sub entity.Subject) (entity.SubjectReview, error) {
	userHash, subHash := user.Hash(), sub.Hash()
	review := entity.SubjectReview{}

	err := checkReviewPermission(DB, userHash, subHash)
	if err != nil {
		return review, err
	}

	snap, err := DB.Restore("users/"+userHash+"/subject_reviews", subHash)
	if err != nil { // user has not reviewed subject
		return review, err
	}

	err = snap.DataTo(&review)
	return review, err
}

func UpdateSubjectReview(DB db.Env, user entity.User, review entity.SubjectReview) error {
	userHash, rvHash := user.Hash(), review.Hash()
	err := checkReviewPermission(DB, userHash, rvHash)
	if err != nil {
		return err
	}

	revRef := DB.Client.Doc("users/" + userHash + "/subject_reviews/" + rvHash)
	subRef := DB.Client.Doc("subjects/" + rvHash)

	err = DB.Client.RunTransaction(DB.Ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		rev, _ := tx.Get(revRef) // get existing review

		if rev.Exists() { // user has already reviewed subject so we must remove it and propagate
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
		err := tx.Set(revRef, review)
		if err != nil {
			return err
		}

		// update subject stats with new review
		for k, v := range review.Review {
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

	return err
}
