package models

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// Professor is an object that represents a USP professor
//
// It is not a DTO and is only used for data collection purposes
type Professor struct {
	CodPes string `firestore:"-"`
	Name   string `firestore:"name"`

	Offerings []Offering `firestore:"-"`
}

// Hash returns SHA256(professor_code)
func (p Professor) Hash() string {
	return utils.SHA256(p.CodPes)
}

// Insert sets a Professor to a given collection. This is usually /institutes/#institute/professors/#professor
func (p Professor) Insert(DB db.Database, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(p.Hash()).Set(DB.Ctx, p)
	return err
}

// Update sets a Professor to a given collection
func (p Professor) Update(DB db.Database, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(p.Hash()).Set(DB.Ctx, p)
	return err
}
