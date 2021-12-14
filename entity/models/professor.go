package models

import "github.com/Projeto-USPY/uspy-backend/db"

// Professor is an object that represents a USP professor
//
// It is not a DTO and is only used for data collection purposes
type Professor struct {
	CodPes string `firestore:"-"`
	Name   string `firestore:"-"`

	Offerings []Offering `firestore:"-"`
}

// Insert is a dummy method for professor. It is necessary because professor must be able to be represented as a db.Writer
func (prof Professor) Insert(DB db.Env, collection string) error {
	return nil
}

// Update is a dummy method for professor. It is necessary because Institute must be able to be represented as a db.Writer
func (prof Professor) Update(DB db.Env, collection string) error {
	return nil
}
