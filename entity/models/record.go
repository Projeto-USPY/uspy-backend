package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// Record is a DTO for a user's subject grade record
//
// It contains some properties that are not mapped to the database and solely used for internal, contextual logic.
type Record struct {
	Subject        string `firestore:"-"`
	Course         string `firestore:"-"`
	Specialization string `firestore:"-"`

	Year     int `firestore:"-"`
	Semester int `firestore:"-"`

	Grade     float64 `firestore:"grade"`
	Status    string  `firestore:"status,omitempty"`
	Frequency int     `firestore:"frequency,omitempty"`
}

// Hash returns SHA256(concat(record_year, record_semester))
//
// This is done in this way so that if a user has more than one record for a given subject, they differ by when the subject was taken.
func (mf Record) Hash() string {
	str := fmt.Sprintf("%d%d", mf.Year, mf.Semester)
	return utils.SHA256(str)
}

// Insert sets a record to a given collection. This is usually /users/#user/final_scores/#final_score/records
func (mf Record) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, mf)
	return err
}

// Update is a dummy method for a record
func (mf Record) Update(DB db.Env, collection string) error { return nil }
