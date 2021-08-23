package models

import (
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

type User struct {
	ID     string `firestore:"-"`
	IDHash string `firestore:"-"` // used only when the ID is unknown

	Name string `firestore:"-"`

	// NameHash and EmailHash are AES encrypted since it has to be decrypted
	NameHash  string `firestore:"name"`
	EmailHash string `firestore:"email"`

	Verified bool `firestore:"verified"` // Email verification
	Banned   bool `firestore:"banned"`

	// bcrypt hashing cause password is more sensitive
	PasswordHash string `firestore:"password"`

	LastUpdate time.Time `firestore:"last_update"`
}

func (u User) Hash() string {
	if u.IDHash != "" {
		return u.IDHash
	}

	return utils.SHA256(u.ID)
}

func NewUser(ID, name, email, password string, lastUpdate time.Time) (*User, error) {
	if nHash, err := utils.AESEncrypt(name, config.Env.AESKey); err != nil {
		return nil, err
	} else {
		if eHash, err := utils.AESEncrypt(email, config.Env.AESKey); err != nil {
			return nil, err
		} else {
			pHash, err := utils.Bcrypt(password)
			if err != nil {
				return nil, err
			}

			return &User{
				ID:           ID,
				Name:         name,
				NameHash:     nHash,
				EmailHash:    eHash,
				PasswordHash: pHash,
				LastUpdate:   lastUpdate,
				Verified:     !config.Env.IsLocal(),
				Banned:       false,
			}, nil
		}
	}
}

func (u User) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(u.Hash()).Set(DB.Ctx, u)
	return err
}

func (u User) Update(DB db.Env, collection string) error {
	updates := make([]firestore.Update, 0)

	if u.PasswordHash != "" {
		updates = append(updates, firestore.Update{
			Path:  "password",
			Value: u.PasswordHash,
		})
	}

	if u.Verified {
		updates = append(updates, firestore.Update{
			Path:  "verified",
			Value: u.Verified,
		})

	}

	_, err := DB.Client.Collection(collection).Doc(u.Hash()).Update(DB.Ctx, updates)
	return err
}
