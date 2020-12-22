package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
)

// Subject describes a subject (example: SMA0356 - Cálculo IV)
type Subject struct {
	Code          string         `form:"code" firestore:"code" binding:"required,alphanum"`
	CourseCode    string         `form:"course" firestore:"course" binding:"required,alphanum"`
	Name          string         `firestore:"name"`
	Description   string         `firestore:"desc"`
	ClassCredits  int            `firestore:"class"`
	AssignCredits int            `firestore:"assign"`
	TotalHours    string         `firestore:"hours"`
	Requirements  []string       `firestore:"requirements"`
	Optional      bool           `firestore:"optional"`
	Stats         map[string]int `firestore:"stats"`
}

func (s Subject) Hash() string {
	str := fmt.Sprintf("%s%s", s.Code, s.CourseCode)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (s Subject) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(s.Hash()).Set(DB.Ctx, s)
	if err != nil {
		return err
	}
	return nil
}
