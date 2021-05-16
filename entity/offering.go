/* package entity contains structs that will be used for backend input validation and DB operations */
package entity

import (
	"crypto/md5"
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
)

// entity.Offering describes an offering of a subject
// Since it is inside a subcollection of a subject, it does not have subject data
type Offering struct {
	Professor string `json:"name" firestore:"professor"`
	CodPes    string `json:"-" firestore:"-"`
	Code      string `json:"code" firestore:"-"`
}

// md5(CodPes)
func (off Offering) Hash() string {
	concat := fmt.Sprint(
		off.CodPes,
	)

	return fmt.Sprintf("%x", md5.Sum([]byte(concat)))
}

func (off Offering) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(off.Hash()).Set(DB.Ctx, off)
	return err
}

func (off Offering) Update(DB db.Env, collection string) error { return nil }
