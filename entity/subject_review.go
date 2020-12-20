package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
)

type SubjectReview struct {
	Subject string          `firestore:"code"`
	Review  map[string]bool `firestore:"scores"`
}

func (sr SubjectReview) Hash() string {
	str := fmt.Sprint(sr.Subject)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (sr SubjectReview) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(sr.Hash()).Set(DB.Ctx, sr)
	if err != nil {
		return err
	}

	return nil
}
