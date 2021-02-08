/* package entity contains structs that will be used for backend input validation and DB operations */
package entity

import (
	"crypto/sha256"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
)

// entity.Grade represents a grade that can be parsed from a user's records PDF
// Example: {10533702, 10.0, 90, "A", "SMA0354", "55041", 1, 2019}
type Grade struct {
	User      string  `json:"-" firestore:"-"`
	Grade     float64 `json:"grade" firestore:"value"`
	Frequency int     `json:"frequency" firestore:"-"`
	Status    string  `json:"status" firestore:"-"`
	Subject   string  `json:"subject" firestore:"-"`
	Course    string  `json:"course" firestore:"-"`
	Semester  int     `json:"semester" firestore:"-"`
	Year      int     `json:"year" firestore:"-"`
}

func (g Grade) Hash() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(g.User)))
}

func (g Grade) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, g)
	if err != nil {
		return err
	}
	return nil
}
