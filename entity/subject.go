package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
)

// Subject describes a subject (example: SMA0356 - Cálculo IV)
type Subject struct {
	Code          string         `json:"code" form:"code" firestore:"code" binding:"required,alphanum"`
	CourseCode    string         `json:"course" form:"course" firestore:"course" binding:"required,alphanum"`
	Name          string         `json:"name" firestore:"name"`
	Description   string         `json:"description" firestore:"desc"`
	ClassCredits  int            `json:"class" firestore:"class"`
	AssignCredits int            `json:"assign" firestore:"assign"`
	TotalHours    string         `json:"hours" firestore:"hours"`
	Requirements  []string       `json:"requirements" firestore:"requirements"`
	Optional      bool           `json:"optional" firestore:"optional"`
	Stats         map[string]int `json:"stats" firestore:"stats"`
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
