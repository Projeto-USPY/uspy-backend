package models

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// Offering is the DTO for an offering of a subject
// Since it is inside a subcollection of a subject, it does not have subject data
//
// It contains some properties that are not mapped to firestore for internal logic and data collection purposes.
type Offering struct {
	CodPes string `firestore:"-"`
	Code   string `firestore:"-"`

	Professor string   `firestore:"professor"`
	Years     []string `firestore:"years"`
}

// Hash returns SHA256(professor_code)
func (off Offering) Hash() string {
	return utils.SHA256(off.CodPes)
}

// Insert sets an offering to a given collection. This is usually /subjects/#subject/offerings
func (off Offering) Insert(DB db.Database, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(off.Hash()).Set(DB.Ctx, off)
	return err
}

// Update sets an offering to a given collection
func (off Offering) Update(DB db.Database, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(off.Hash()).Set(DB.Ctx, off)
	return err
}
