package models

import (
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// User is the DTO for a registered user
//
// It also contains non mapped properties used for internal contextual logic.
type User struct {
	ID     string `firestore:"-"`
	IDHash string `firestore:"-"` // used only when the ID is unknown

	Name string `firestore:"-"`

	// NameHash is AES encrypted since it has to be decrypted
	NameHash string `firestore:"name"`

	// EmailHash is SHA256 encrypted because it must be queried in signup
	EmailHash string `firestore:"email"`

	Verified bool `firestore:"verified"` // Email verification
	Banned   bool `firestore:"banned"`

	// bcrypt hashing cause password is more sensitive
	PasswordHash string `firestore:"password"`

	LastUpdate      time.Time        `firestore:"last_update"`
	TranscriptYears map[string][]int `firestore:"transcript_years"` // map of years as keys and array of semesters as values
}

// Hash returns SHA256(user_id)
func (u User) Hash() string {
	if u.IDHash != "" {
		return u.IDHash
	}

	return utils.SHA256(u.ID)
}

// NewUser creates a new user. It takes raw data and processes all the encrypted data
//
// User email verification is also bypassed in dev or local environments
func NewUser(ID, name string, lastUpdate time.Time, transcriptYears map[string][]int) (*User, error) {
	nHash, err := utils.AESEncrypt(name, config.Env.AESKey)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:              ID,
		Name:            name,
		NameHash:        nHash,
		LastUpdate:      lastUpdate,
		TranscriptYears: transcriptYears,
		Verified:        config.Env.IsLocal() || config.Env.IsDev(),
		Banned:          false,
	}, nil

}

// Insert sets a user object to a given collection. This is usually /users
func (u User) Insert(DB db.Database, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(u.Hash()).Set(DB.Ctx, u)
	return err
}

// CompleteSignup sets a user's email and password and verifies an user. Collection is usually /users
func CompleteSignup(DB db.Database, userHash, collection, email, password string) error {
	eHash := utils.SHA256(email)
	pHash, err := utils.Bcrypt(password)
	if err != nil {
		return err
	}

	_, err = DB.Client.Collection(collection).Doc(userHash).Update(DB.Ctx, []firestore.Update{
		{
			Path:  "email",
			Value: eHash,
		},
		{
			Path:  "password",
			Value: pHash,
		},
		{
			Path:  "verified",
			Value: true,
		},
	})

	return err
}

// Update sets a user object to a given collection. This is usually /users
//
// This method only allows updating the password
// TODO: Use MergeWithout to specifically mention non-updatable fields
func (u User) Update(DB db.Database, collection string) error {
	updates := make([]firestore.Update, 0)

	if u.PasswordHash != "" {
		updates = append(updates, firestore.Update{
			Path:  "password",
			Value: u.PasswordHash,
		})
	}

	_, err := DB.Client.Collection(collection).Doc(u.Hash()).Update(DB.Ctx, updates)
	return err
}
