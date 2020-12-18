package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"google.golang.org/api/iterator"
)

// Professor represents a ICMC professor (example: {Moacir Ponti SCC})
type Professor struct {
	CodPes     int    `firestore:"code,omitempty"`
	Name       string `firestore:"name,omitempty"`
	Department string `firestore:"dep,omitempty"`

	Stats     map[string]int `firestore:"stats,omitempty"`
	Offerings []string       `firestore:"offeringsIDs,omitempty"`
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
