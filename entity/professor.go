/* package entity contains structs that will be used for backend input validation and DB operations */
package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/db"
	"google.golang.org/api/iterator"
)

// entity.Professor represents a ICMC professor
// Example: {1234567, "Fulano da Silva", SCC, Stats, Offerings}
/* Stats will be a map of review statistics, that looks like:
stats = {
	"total": 5,
	"worth_it": 3,
	"is_difficult": 4,
	"has_attendance": 0,
	"has_didactics": 4
}
*/
// Offerings contains a list of all offerings the professor was responsible for since 2010
type Professor struct {
	CodPes     int    `json:"code" firestore:"code"`
	Name       string `json:"name" firestore:"name"`
	Department string `json:"department" firestore:"dep,omitempty"`

	Stats     map[string]int `json:"stats" firestore:"stats"`
	Offerings []string       `json:"offerings" firestore:"offerings_ids"`
}

func (prof Professor) WithOfferings(DB db.Env) (Professor, error) {
	offs, err := prof.GetProfessorOfferingIDs(DB)
	if err != nil {
		return Professor{}, err
	}

	prof.Offerings = offs
	return prof, nil
}

func (prof Professor) GetProfessorOfferingIDs(DB db.Env) ([]string, error) {
	col := DB.Client.Collection("offerings")
	offs := col.Where("data.professor", "==", prof.CodPes)
	iter := offs.Documents(DB.Ctx)
	defer iter.Stop()

	offeringIDs := make([]string, 0, 500)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		offeringIDs = append(offeringIDs, doc.Ref.ID)
	}

	return offeringIDs, nil
}

// md5(concat(professor.CodPes))
func (prof Professor) Hash() string {
	str := fmt.Sprint(prof.CodPes)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (prof Professor) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(prof.Hash()).Set(DB.Ctx, prof)
	if err != nil {
		return err
	}

	return nil
}
