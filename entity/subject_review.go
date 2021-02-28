/* Package db contains useful functions related to the Firestore Database */
package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"reflect"
	"sort"
)

func validateSubjectReview(f1 validator.FieldLevel) bool {
	keys := f1.Field().MapKeys()
	keysStr := make([]string, 0)

	for _, k := range keys {
		keysStr = append(keysStr, k.String())
	}

	categories := getCategories()

	if len(keys) != len(categories) {
		return false
	}

	sort.Strings(keysStr)
	sort.Strings(categories)

	for i := 0; i < len(keys); i++ {
		if !reflect.DeepEqual(keys[i].String(), categories[i]) {
			return false
		}
	}

	return true
}

// entity.SubjectReview represents a review made to a subject by a user
// Example: {"SMA0354", "55041", map[string]interface{}{"worth_it": true}}
type SubjectReview struct {
	Subject        string                 `json:"-" firestore:"-" binding:"required,alphanum"`
	Course         string                 `json:"-" firestore:"-" binding:"required,alphanum"`
	Specialization string                 `json:"-" firestore:"-" binding:"required,alphanum"`
	Review         map[string]interface{} `json:"categories" firestore:"categories" binding:"required,validateSubjectReview"`
}

func (sr SubjectReview) Hash() string {
	str := fmt.Sprintf("%s%s%s", sr.Subject, sr.Course, sr.Specialization)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (sr SubjectReview) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(sr.Hash()).Set(DB.Ctx, sr)
	if err != nil {
		return err
	}

	return nil
}

func getCategories() []string {
	return []string{"worth_it"}
}
