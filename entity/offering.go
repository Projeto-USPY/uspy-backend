/* package entity contains structs that will be used for backend input validation and DB operations */
package entity

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// entity.Offering describes an offering of a subject
// Since it is inside a subcollection of a subject, it does not have subject data
type Offering struct {
	Professor string `json:"name" firestore:"professor"`
	CodPes    string `json:"-" firestore:"-"`
	Code      string `json:"code" firestore:"-"`
	Year      string `json:"year" firestore:"year"`
}

// sha256(CodPes)
func (off Offering) Hash() string {
	return utils.SHA256(off.CodPes)
}

func (off Offering) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(off.Hash()).Set(DB.Ctx, off)
	return err
}

func (off Offering) Update(DB db.Env, collection string) error { return nil }
