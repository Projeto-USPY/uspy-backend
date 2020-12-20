package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
)

// Subject describes a subject (example: SMA0356 - Cálculo IV)
type Subject struct {
	Code          string         `firestore:"code,omitempty"`
	Name          string         `firestore:"name,omitempty"`
	Description   string         `firestore:"desc,omitempty"`
	ClassCredits  int            `firestore:"class,omitempty"`
	AssignCredits int            `firestore:"assign,omitempty"`
	TotalHours    string         `firestore:"hours,omitempty"`
	Requirements  []string       `firestore:"requirements,omitempty"`
	Optional      bool           `firestore:"optional,omitempty"`
	Stats         map[string]int `firestore:"stats,omitempty"`
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
