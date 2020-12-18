package entity

import (
	"cloud.google.com/go/firestore"
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
)

// Subject describes a subject (example: SMA0356 - Cálculo IV)
type Subject struct {
	Code          string                   `firestore:"code"`
	Name          string                   `firestore:"name"`
	Description   string                   `firestore:"desc"`
	ClassCredits  int                      `firestore:"class"`
	AssignCredits int                      `firestore:"assign"`
	TotalHours    string                   `firestore:"hours"`
	Requirements  []string                 `firestore:"requirements"`
	Optional      bool                     `firestore:"optional"`
	Grades        *firestore.CollectionRef `firestore:"grades"`
}

func (s Subject) Hash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s.Code)))
}

func (s Subject) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(s.Hash()).Set(DB.Ctx, s)
	if err != nil {
		return err
	}
	return nil
}
