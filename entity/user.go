package entity

import (
	"crypto/sha256"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"time"
)

type User struct {
	Login            string    `json:"login" firestore:"login,omitempty" binding:"required"`
	LastRegistration time.Time `firestore:"lastRegistration,serverTimestamp"`

	// Password is only used because of REST requests, do not store in DB!!!
	Password     string `json:"pwd" firestore:"-" binding:"required"`
	PasswordHash string `firestore:"password,omitempty"` // use strong Hashing!
}

// sha256 cause user data is more sensitive
func (u User) Hash() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(u.Login)))
}

func (u User) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(u.Hash()).Set(DB.Ctx, u)
	if err != nil {
		return err
	}

	return nil
}
