package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

type Grade struct {
	User           string `firestore:"-"`
	Subject        string `firestore:"-"`
	Frequency      int    `firestore:"-"`
	Status         string `firestore:"-"`
	Course         string `firestore:"-"`
	Specialization string `firestore:"-"`
	Semester       int    `firestore:"-"`
	Year           int    `firestore:"-"`

	Value float64 `firestore:"value"`
}

func (g Grade) Hash() string {
	str := fmt.Sprintf("%x%x", g.Year, g.Semester)
	return utils.SHA256(str)
}

func (g Grade) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, g)
	return err
}

func (g Grade) Update(DB db.Env, collection string) error { return nil }
