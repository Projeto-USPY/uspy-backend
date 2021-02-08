// package account contains functions that implement backend-db communication for every /account endpoint
package account

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/iddigital"
	"golang.org/x/crypto/bcrypt"
)

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
func Login(DB db.Env, u entity.User) error {
	snap, err := DB.Restore("users", u.Hash())
	if err != nil {
		return err
	}

	var storedUser entity.User
	err = snap.DataTo(&storedUser)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(u.Password))
}
