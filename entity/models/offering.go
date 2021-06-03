package models

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

type Offering struct {
	CodPes string
	Code   string

	Professor string `firestore:"professor"`
	Year      string `firestore:"year"`
}

func (off Offering) Hash() string {
	return utils.SHA256(off.CodPes)
}

func (off Offering) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(off.Hash()).Set(DB.Ctx, off)
	return err
}

func (off Offering) Update(DB db.Env, collection string) error { return nil }
