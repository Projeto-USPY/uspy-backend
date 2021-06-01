/* package entity contains structs that will be used for backend input validation and DB operations */
package entity

import (
	"fmt"
	"reflect"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// entity.Course represents a course/major
// Example: {"Bacharelado em Ciências de Computação", "55041", []Subjects{...}, map[string]string{"SMA0356": "Cálculo IV", ...}}
type Course struct {
	Name           string            `json:"name" firestore:"name"`
	Code           string            `json:"code" firestore:"code"`
	Specialization string            `json:"specialization" firestore:"specialization"`
	Subjects       []Subject         `json:"-" firestore:"-"`
	SubjectCodes   map[string]string `json:"subjects" firestore:"subjects"`
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
	updates := make([]firestore.Update, 0)
	fields := reflect.TypeOf(c)
	values := reflect.ValueOf(c)

	for i := 0; i < fields.NumField(); i++ {
		fieldValue := values.Field(i).Interface()
		if tag := fields.Field(i).Tag.Get("firestore"); tag != "-" {
			updates = append(updates, firestore.Update{Path: tag, Value: fieldValue})
		}
	}

	_, err := DB.Client.Collection(collection).Doc(c.Hash()).Update(DB.Ctx, updates)
	return err
}
