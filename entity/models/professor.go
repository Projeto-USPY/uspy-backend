package models

import "github.com/Projeto-USPY/uspy-backend/db"

type Professor struct {
	CodPes string `firestore:"-"`
	Name   string `firestore:"-"`

	Offerings []Offering `firestore:"-"`
}

func (prof Professor) Insert(DB db.Env, collection string) error {
	return nil
}

func (prof Professor) Update(DB db.Env, collection string) error {
	return nil
}
