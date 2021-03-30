/* package entity contains structs that will be used for backend input validation and DB operations */
package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/db"
)

// entity.ProfessorReview represents a review made by an user to a specific professor\
// Example: {12345678, Review}
/* Review will be a map that may look like:
review = {
	"worth_it": true,
	"is_difficult": true,
	"has_attendance": false,
	"has_didactics": true,
}
*/
type ProfessorReview struct {
	CodPes int             `json:"code" firestore:"code"`
	Review map[string]bool `json:"reviews" firestore:"reviews"`
}

func (pr ProfessorReview) Hash() string {
	str := fmt.Sprint(pr.CodPes)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (pr ProfessorReview) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(pr.Hash()).Set(DB.Ctx, pr)
	if err != nil {
		return err
	}

	return nil
}
