/* Package db contains useful functions related to the Firestore Database */
package entity

import (
	"crypto/md5"
	"fmt"

	"github.com/tpreischadt/ProjetoJupiter/db"
)

// entity.Subject describes a subject
// Example: {"SMA0354", "55041", "Cálculo IV", "...", 4, 1, 60, []string{"Cálculo I", ...}, false, Stats}
// Stats will be a map of review stats that may look like:
/*
	stats = {
		"total": 3,
		"worth_it": 2
	}
*/
type Subject struct {
	Code             string                   `json:"code" form:"code" firestore:"code" binding:"required,alphanum"`
	CourseCode       string                   `json:"course" form:"course" firestore:"course" binding:"required,alphanum"`
	Specialization   string                   `json:"specialization" form:"specialization" firestore:"specialization" binding:"required,alphanum"`
	Name             string                   `json:"name" firestore:"name"`
	Description      string                   `json:"description" firestore:"desc"`
	Semester         int                      `json:"semester" firestore:"semester"`
	ClassCredits     int                      `json:"class" firestore:"class"`
	AssignCredits    int                      `json:"assign" firestore:"assign"`
	TotalHours       string                   `json:"hours" firestore:"hours"`
	Requirements     map[string][]Requirement `json:"requirements" firestore:"requirements"`
	TrueRequirements []Requirement            `json:"-" firestore:"true_requirements"`
	Optional         bool                     `json:"optional" firestore:"optional"`
	Stats            map[string]int           `json:"stats" firestore:"stats"`
}

func (s Subject) Hash() string {
	str := fmt.Sprintf("%s%s%s", s.Code, s.CourseCode, s.Specialization)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (s Subject) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(s.Hash()).Set(DB.Ctx, s)
	if err != nil {
		return err
	}
	return nil
}
