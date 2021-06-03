package models

import (
	"github.com/Projeto-USPY/uspy-backend/db"
)

type Institute struct {
	Name    string   `firestore:"-"`
	Code    string   `firestore:"-"`
	Courses []Course `firestore:"-"`

	Professors []Professor `firestore:"-"`
}

func (i Institute) Insert(DB db.Env, collection string) error { return nil }

func (i Institute) Update(DB db.Env, collection string) error { return nil }
