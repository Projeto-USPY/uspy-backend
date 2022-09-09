package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// Major is the DTO for a course/major. It represents a tuple (course, specialization)
//
// It is used for storing which courses a user has records of
type Major struct {
	Code           string `firestore:"course"`
	Specialization string `firestore:"specialization"`
}

// NewMajorFromController is a constructor. It takes a major controller and returns a model.
func NewMajorFromController(major *controllers.Major) *Major {
	return &Major{
		Code:           major.Course,
		Specialization: major.Specialization,
	}
}

// Hash returns SHA256(concat(course, specialization))
func (m Major) Hash() string {
	str := fmt.Sprintf("%s%s", m.Code, m.Specialization)
	return utils.SHA256(str)
}

// Insert sets a major to a given collection. This is usually /users/#user/majors
func (m Major) Insert(DB db.Database, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(m.Hash()).Set(DB.Ctx, m)
	return err
}

// Update is a dummy method for a major
func (m Major) Update(DB db.Database, collection string) error { return nil }
