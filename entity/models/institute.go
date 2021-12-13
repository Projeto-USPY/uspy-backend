package models

import (
	"github.com/Projeto-USPY/uspy-backend/db"
)

// Institute represents an institute or collection of courses and professors.
//
// It is not a DTO and therefore not mapped to any entity in the database
// This entity is used for data collection only.
type Institute struct {
	Name    string   `firestore:"-"`
	Code    string   `firestore:"-"`
	Courses []Course `firestore:"-"`

	Professors []Professor `firestore:"-"`
}

// Insert is a dummy method for Institute. It is necessary because Institute must be able to be represented as a db.Writer
func (i Institute) Insert(DB db.Env, collection string) error { return nil }

// Update is a dummy method for Institute. It is necessary because Institute must be able to be represented as a db.Writer
func (i Institute) Update(DB db.Env, collection string) error { return nil }
