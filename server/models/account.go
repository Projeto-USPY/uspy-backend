package models

import (
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/iddigital"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

// GenerateJWT generates a JWT from user struct
func GenerateJWT(user entity.User) (jwtString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":      user.Login,
		"timestamp": time.Now().Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	jwtString, err = token.SignedString([]byte(secret))

	return
}

// ValidateJWT takes a JWT token string and validates it
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if token == nil || !token.Valid {
		return nil, err
	}

	return token, nil
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

		// store all user records
		objs = append(objs, db.Object{
			Collection: "users/" + u.Hash() + "/final_scores",
			Data:       mf,
		})

		// add grade to "global" grades collection
		gradeObj := entity.Grade{
			User:  u.Login,
			Grade: g.Grade,
		}

		objs = append(objs, db.Object{
			Collection: "subjects/" + entity.Subject{Code: g.Subject, CourseCode: g.Course}.Hash() + "/grades",
			Data:       gradeObj,
		})
	}

	writeErr := DB.BatchWrite(objs)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

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
