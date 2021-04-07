/* package entity contains structs that will be used for backend input validation and DB operations */
package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/db"
)

// entity.FinalScore is a user's final score that is stored in the Firestore DB
// Example: {10.0, "A", 2019, 1, 90, offeringHash}
// Be aware the final score does not include the entity.Subject data, because that is included in the offeringHash
type FinalScore struct {
	Grade        float64 `json:"grade" firestore:"grade"`
	Status       string  `json:"status" firestore:"status"`
	Year         int     `json:"-" firestore:"-"`
	Semester     int     `json:"-" firestore:"-"`
	Frequency    int     `json:"frequency" firestore:"frequency"`
	OfferingHash string  `json:"-" firestore:"offering,omitempty"`
}

func (mf FinalScore) Hash() string {
	str := fmt.Sprintf("%d%d", mf.Year, mf.Semester)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (mf FinalScore) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, mf)
	return err
}
