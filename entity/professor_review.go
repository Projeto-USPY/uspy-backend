package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
)

type ProfessorReview struct {
	CodPes int            `firestore:"code,omitempty"`
	Review map[string]int `firestore:"scores,omitempty"`
}

func (pr ProfessorReview) Hash() string {
	str := fmt.Sprint(pr.CodPes)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (pr ProfessorReview) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(pr.Hash()).Set(DB.Ctx, pr)
	if err != nil {
		return err
	}

	return nil
}
