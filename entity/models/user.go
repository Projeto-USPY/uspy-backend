package models

import (
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

type User struct {
	ID   string `firestore:"-"`
	Name string `firestore:"-"`

	// NameHash is AES encrypted since it has to be decrypted
	NameHash string `firestore:"name"`
	// bcrypt hashing cause password is more sensitive
	PasswordHash string `firestore:"password"`

	LastUpdate time.Time `firestore:"last_update"`
}

func (u User) Hash() string {
	return utils.SHA256(u.ID)
}

func NewUser(ID, name, password string, lastUpdate time.Time) (*User, error) {
	nHash, err := utils.AESEncrypt(name, config.Env.AESKey)
	if err != nil {
		return nil, err
	}

	pHash, err := utils.Bcrypt(password)
	if err != nil {
		return nil, err
	}

	return &User{ID: ID, Name: name, NameHash: nHash, PasswordHash: pHash, LastUpdate: lastUpdate}, nil
}

func (u User) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(u.Hash()).Set(DB.Ctx, u)
	return err
}

func (u User) Update(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(u.Hash()).Update(DB.Ctx, []firestore.Update{{
		Path:  "password",
		Value: u.PasswordHash,
	}})
	return err
}
