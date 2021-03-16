/* package entity contains structs that will be used for backend input validation and DB operations */
package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/db"
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
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (c Course) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(c.Hash()).Set(DB.Ctx, c)
	if err != nil {
		return err
	}

	return nil
}
