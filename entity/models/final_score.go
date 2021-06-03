package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

type FinalScore struct {
	Year     int
	Semester int

	Grade     float64 `firestore:"grade"`
	Status    string  `firestore:"status"`
	Frequency int     `firestore:"frequency"`
}

func (mf FinalScore) Hash() string {
	str := fmt.Sprintf("%d%d", mf.Year, mf.Semester)
	return utils.SHA256(str)
}

func (mf FinalScore) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, mf)
	return err
}

func (mf FinalScore) Update(DB db.Env, collection string) error { return nil }
