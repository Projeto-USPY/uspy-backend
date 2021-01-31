package entity

import (
	"crypto/sha256"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	// only used because of REST requests, do not store in DB!!!
	Login    string `json:"login" firestore:"-" binding:"required"`
	Password string `json:"pwd" firestore:"-" binding:"required"`

	// bcrypt hashing cause password is more sensitive
	PasswordHash string `firestore:"password"`

	LastUpdate time.Time `firestore:"lastUpdate"`
}

func HashPassword(str string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", pass), nil
}

func (u User) WithHash() (User, error) {
	pHash, err := HashPassword(u.Password)
	if err != nil {
		return User{}, err
	}
	u.PasswordHash = pHash
	return u, nil
}

func (u User) Hash() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(u.Login)))
}

func (u User) Insert(DB db.Env, collection string) error {
	u, err := u.WithHash()
	if err != nil {
		return err
	}
	_, err = DB.Client.Collection(collection).Doc(u.Hash()).Set(DB.Ctx, u)
	if err != nil {
		return err
	}

	return nil
}

func (u User) Signup(DB db.Env, recs Records) error {
	objs := []db.Object{
		db.Object{
			Collection: "users",
			Doc:        u.Hash(),
			Data:       u,
		},
	}

	for _, g := range recs.Grades {
		mf := FinalScore{
			Grade:     g.Grade,
			Status:    g.Status,
			Frequency: g.Frequency,
		}

		// store all user records
		objs = append(objs, db.Object{
			Collection: "users/" + u.Hash() + "/final_scores",
			Data:       mf,
		})

		// add grade to "global" grades collection
		gradeObj := Grade{
			User:  u.Login,
			Grade: g.Grade,
		}

		objs = append(objs, db.Object{
			Collection: "subjects/" + Subject{Code: g.Subject, CourseCode: g.Course}.Hash() + "/grades",
			Data:       gradeObj,
		})
	}

	writeErr := DB.BatchWrite(objs)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
