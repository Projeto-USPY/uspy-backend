// package account contains functions that implement backend-db communication for every /account endpoint
package account

import (
	"fmt"
	"reflect"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/iddigital"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// Profile retrieves the user profile. In v1 that only contains name and id info
func Profile(DB db.Env, u entity.User) (entity.User, error) {
	var storedUser entity.User
	snap, err := DB.Restore("users", u.Hash())
	if err != nil {
		return entity.User{}, err
	}
	err = snap.DataTo(&storedUser)
	if err != nil {
		return entity.User{}, err
	}
	return storedUser, nil
}

// Signup inserts a new user into the DB
func Signup(DB db.Env, u entity.User, recs iddigital.Records) error {
	objs := []db.Object{
		{
			Collection: "users",
			Doc:        u.Hash(),
			Data:       u,
		},
	}

	for _, g := range recs.Grades {
		mf := entity.FinalScore{
			Grade:     g.Grade,
			Status:    g.Status,
			Frequency: g.Frequency,
		}

		subHash := entity.Subject{Code: g.Subject, CourseCode: g.Course}.Hash()

		// store all user records
		objs = append(objs, db.Object{
			Collection: "users/" + u.Hash() + "/final_scores/" + subHash + "/records",
			Doc:        mf.Hash(),
			Data:       mf,
		})

		// add grade to "global" grades collection
		gradeObj := entity.Grade{
			User:  u.Login,
			Grade: g.Grade,
		}

		objs = append(objs, db.Object{
			Collection: "subjects/" + subHash + "/grades",
			Data:       gradeObj,
		})
	}

	// write atomically
	writeErr := DB.BatchWrite(objs)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

// Login performs the user login by comparing the inserted password and the stored hash
func Login(DB db.Env, u entity.User) (entity.User, error) {
	snap, err := DB.Restore("users", u.Hash())
	if err != nil {
		return entity.User{}, err
	}

	var storedUser entity.User
	err = snap.DataTo(&storedUser)
	if err != nil {
		return entity.User{}, err
	}

	return storedUser, bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(u.Password))
}

// ChangePassword changes the current password hash to a new one
func ChangePassword(DB db.Env, u entity.User, newPassword string) error {
	newHash, err := entity.HashPassword(newPassword)
	if err != nil {
		return err
	}
	pwdUpdates := []firestore.Update{{Path: "password", Value: newHash}}
	return DB.Update(u.Hash(), "users", pwdUpdates)
}

// Remove removes the user from Firestore, and all the data associated with it
func Remove(DB db.Env, u entity.User) error {
	userHash := u.Hash()

	userRef := DB.Client.Doc("users/" + userHash)
	finalScoresDocs, err := userRef.Collection("final_scores").DocumentRefs(DB.Ctx).GetAll()
	if err != nil {
		return err
	}

	subjectReviewsDocs, err := userRef.Collection("subject_reviews").DocumentRefs(DB.Ctx).GetAll()
	if err != nil {
		return err
	}

	err = DB.Client.RunTransaction(DB.Ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Do all the reading first

		var wg sync.WaitGroup
		channelErr := make(chan error, len(finalScoresDocs)*100)

		for _, subRef := range finalScoresDocs {
			go func() {
				wg.Add(1)
				defer wg.Done()

				recordsDocs, err := subRef.Collection("records").DocumentRefs(DB.Ctx).GetAll()
				if err != nil {
					channelErr <- err
				}

				gradesCol := DB.Client.Collection("subjects/" + subRef.ID + "/grades")

				// get grades to remove
				for _, recordRef := range recordsDocs {
					go func() {
						wg.Add(1)
						defer wg.Done()

						// get value of record
						snap, err := recordRef.Get(ctx)
						if err != nil {
							channelErr <- err
						}

						// read final score
						var score entity.FinalScore
						if err = snap.DataTo(&score); err != nil {
							channelErr <- err
						}

						// finds one grade in subject/subject_id/grades where the grade is the same as score.Grade
						gradeDocsToRemove, err := gradesCol.Where("value", "==", score.Grade).Limit(1).Documents(ctx).GetAll()
						if err != nil {
							channelErr <- err
						}

						// delete the grade documents (there must be exactly one)
						for _, gradeSnap := range gradeDocsToRemove {
							_, err := gradeSnap.Ref.Delete(ctx)
							if err != nil {
								channelErr <- err
							}
						}
					}()
				}
			}()
		}

		wg.Wait()
		close(channelErr)
		for e := range channelErr {
			if e != nil {
				return e
			}
		}

		// For all review documents
		for _, reviewRef := range subjectReviewsDocs {
			rev, err := reviewRef.Get(ctx) // get review snapshot
			if err != nil {
				return err
			}

			categories, err := rev.DataAt("categories") // get existing review categories
			if err != nil {
				return err
			}

			subRef := DB.Client.Doc("subjects/" + reviewRef.ID)

			// For all of the categories
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

		userRef.Delete(ctx) // deletes the user

		return nil
	})

	return err
}
