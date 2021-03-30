/* Package db contains useful functions related to the Firestore Database */
package entity

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

// entity.User represents a user not only in the Firestore DB but also as input for some endpoints
// Example: {"12345678", "mypwd123", true, hash("mypwd123"), LastUpdate}
// LastUpdate is the creation time of the last PDF records submitted by the user, either on signup or profile update
type User struct {
	// only used because of REST requests, do not store in DB!!!
	Login    string `json:"login" firestore:"-" binding:"required,numeric"`
	Password string `json:"pwd" firestore:"-" binding:"required"`
	Remember bool   `json:"remember" firestore:"-"`

	// bcrypt hashing cause password is more sensitive
	PasswordHash string `firestore:"password"`

	// Name is sensitive, do not store in DB!!!
	Name string `json:"name" firestore:"-"`

	// NameHash is AES encrypted since it has to be decrypted
	NameHash string `firestore:"name"`

	LastUpdate time.Time `firestore:"last_update"`
}

func HashPassword(str string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(pass), nil
}

type UserOptions interface {
	Apply(*User) error
}

type WithPasswordHash struct{}

func (WithPasswordHash) Apply(u *User) error {
	pHash, err := HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.PasswordHash = pHash
	return nil
}

type WithNameHash struct{}

func (WithNameHash) Apply(u *User) error {
	if key, ok := os.LookupEnv("AES_KEY"); ok {
		nHash, err := utils.AESEncrypt(u.Name, key)
		if err != nil {
			return err
		}
		u.NameHash = nHash
	} else {
		return errors.New("AES_KEY 128/196/256-bit key env variable was not provided")
	}

	return nil
}

func NewUserWithOptions(login, password, name string, lastUpdate time.Time, opts ...UserOptions) (User, error) {
	u := User{Login: login, Password: password, Name: name, LastUpdate: lastUpdate}
	for _, opt := range opts {
		if err := opt.Apply(&u); err != nil {
			return User{}, err
		}
	}

	return u, nil
}

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
