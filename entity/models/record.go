package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

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

func (mf Record) Hash() string {
	str := fmt.Sprintf("%d%d", mf.Year, mf.Semester)
	return utils.SHA256(str)
}

func (mf Record) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, mf)
	return err
}

func (mf Record) Update(DB db.Env, collection string) error { return nil }
