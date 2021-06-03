package models

import (
	"fmt"
	"reflect"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

type Subject struct {
	Code             string                   `firestore:"code"`
	CourseCode       string                   `firestore:"course"`
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

func (s Subject) Hash() string {
	str := fmt.Sprintf("%s%s%s", s.Code, s.CourseCode, s.Specialization)
	return utils.SHA256(str)
}

func NewSubjectFromController(sub *controllers.Subject) *Subject {
	return &Subject{Code: sub.Code, CourseCode: sub.CourseCode, Specialization: sub.Specialization}
}

func (s Subject) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(s.Hash()).Set(DB.Ctx, s)
	return err
}

func (s Subject) Update(DB db.Env, collection string) error {
	updates := make([]firestore.Update, 0)
	fields := reflect.TypeOf(s)
	values := reflect.ValueOf(s)

	for i := 0; i < fields.NumField(); i++ {
		fieldValue := values.Field(i).Interface()
		if tag := fields.Field(i).Tag.Get("firestore"); tag != "-" && tag != "stats" {
			updates = append(updates, firestore.Update{Path: tag, Value: fieldValue})
		}
	}

	_, err := DB.Client.Collection(collection).Doc(s.Hash()).Update(DB.Ctx, updates)
	return err
}
