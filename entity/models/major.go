package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

type Major struct {
	Course         string `firestore:"course"`
	Specialization string `firestore:"specialization"`
}

func (m Major) Hash() string {
	str := fmt.Sprintf("%s%s", m.Course, m.Specialization)
	return utils.SHA256(str)
}

func (m Major) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(m.Hash()).Set(DB.Ctx, m)
	return err
}

func (m Major) Update(DB db.Env, collection string) error { return nil }
