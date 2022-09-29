package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	firestoreUtils "github.com/Projeto-USPY/uspy-backend/entity/models/utils"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// Subject is the DTO for a subject.
//
// It is probably one of the most used objects.
type Subject struct {
	Code             string                   `firestore:"code"`
	CourseCode       string                   `firestore:"course"`
	CourseName       string                   `firestore:"course_name"`
	Specialization   string                   `firestore:"specialization"`
	Name             string                   `firestore:"name"`
	Description      string                   `firestore:"desc"`
	Semester         int                      `firestore:"semester"`
	ClassCredits     int                      `firestore:"class"`
	AssignCredits    int                      `firestore:"assign"`
	TotalHours       string                   `firestore:"hours"`
	Requirements     map[string][]Requirement `firestore:"requirements"`
	TrueRequirements []Requirement            `firestore:"true_requirements"`
	Optional         bool                     `firestore:"optional"`
	Stats            map[string]int           `firestore:"stats"`
}

// Hash returns SHA256(concat(code, course, specialization))
func (s Subject) Hash() string {
	str := fmt.Sprintf("%s%s%s", s.Code, s.CourseCode, s.Specialization)
	return utils.SHA256(str)
}

// NewSubjectFromController is a constructor. It takes a subject controller and returns a model.
func NewSubjectFromController(sub *controllers.Subject) *Subject {
	return &Subject{Code: sub.Code, CourseCode: sub.CourseCode, Specialization: sub.Specialization}
}

// Insert sets a subject to a given collection. This is usually /subjects
func (s Subject) Insert(DB db.Database, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(s.Hash()).Set(DB.Ctx, s)
	return err
}

// Update sets a subject to a given collection. This is usually /subjects
//
// This method prohibits from changing the stats map
func (s Subject) Update(DB db.Database, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(s.Hash()).Set(DB.Ctx, s, firestoreUtils.MergeWithout(
		s,
		"stats",
	))
	return err
}
