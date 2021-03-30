// package account contains functions that implement backend-db communication for every /account endpoint
package account

import (
	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
	"github.com/Projeto-USPY/uspy-backend/iddigital"
	"golang.org/x/crypto/bcrypt"
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

		subHash := entity.Subject{Code: g.Subject, CourseCode: g.Course, Specialization: g.Specialization}.Hash()

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
