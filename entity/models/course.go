package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// Course is the DTO for a course
// Example: {"Bacharelado em Ciências de Computação", "55041", []Subjects{...}, map[string]string{"SMA0356": "Cálculo IV", ...}}
type Course struct {
	Name           string            `firestore:"name"`
	Code           string            `firestore:"code"`
	Specialization string            `firestore:"specialization"`
	Shift          string            `firestore:"shift"`
	SubjectCodes   map[string]string `firestore:"subjects"`

	Subjects []Subject `firestore:"-"`
}

// Hash returns SHA256(concat(course, specialization))
func (c Course) Hash() string {
	str := fmt.Sprintf("%s%s", c.Code, c.Specialization)
	return utils.SHA256(str)
}

// Insert sets a course to a given collection. This is usually /institutes/#institute/courses
func (c Course) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(c.Hash()).Set(DB.Ctx, c)
	return err
}

// Update sets a course to a given collection
func (c Course) Update(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(c.Hash()).Set(DB.Ctx, c)
	return err
}
