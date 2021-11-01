package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// entity.Course represents a course/major
// Example: {"Bacharelado em Ciências de Computação", "55041", []Subjects{...}, map[string]string{"SMA0356": "Cálculo IV", ...}}
type Course struct {
	Name           string            `firestore:"name"`
	Code           string            `firestore:"code"`
	Specialization string            `firestore:"specialization"`
	SubjectCodes   map[string]string `firestore:"subjects"`

	Subjects []Subject `firestore:"-"`
}

func (c Course) Hash() string {
	str := fmt.Sprintf("%s%s", c.Code, c.Specialization)
	return utils.SHA256(str)
}

func (c Course) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(c.Hash()).Set(DB.Ctx, c)
	return err
}

func (c Course) Update(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(c.Hash()).Set(DB.Ctx, c)
	return err
}
