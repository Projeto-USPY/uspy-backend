package db

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"google.golang.org/api/iterator"
)

type ProfessorDB struct {
	HashID    string           // md5(concat(professor.CodPes))
	Offerings []string         `firestore:"offeringsIDs"`
	Stats     map[string]int   `firestore:"stats"`
	Professor entity.Professor `firestore:"data"`
}

// NewProfessor constructor
func NewProfessor(ent entity.Professor, DB Env) (*ProfessorDB, error) {
	col := DB.Client.Collection("offerings")
	offs := col.Where("data.professor", "==", ent.CodPes)
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

	prof := ProfessorDB{
		Offerings: offeringIDs,
		Stats: map[string]int{
			"sumDidactics": 0,
			"sumRigorous":  0,
		},
		Professor: ent,
	}

	prof.Hash()
	return &prof, nil
}

/* NewProfessorWithOfferings is the same as NewProfessor but you must inform the offerings.
Use this method instead of NewProfessor in order to reduce firestore reads when inserting professors. */
func NewProfessorWithOfferings(ent entity.Professor, offs []entity.Offering) (*ProfessorDB, error) {
	offHashes := make([]string, 0, 500)
	for _, off := range offs {
		if off.Professor != ent.CodPes {
			return nil, errors.New("invalid offering")
		} else {
			offDB := NewOffering(off)
			offHashes = append(offHashes, offDB.HashID)
		}
	}
	prof := ProfessorDB{
		Offerings: offHashes,
		Stats: map[string]int{
			"sumDidactics": 0,
			"sumRigorous":  0,
		},
		Professor: ent,
	}
	prof.HashID = prof.Hash()
	return &prof, nil
}

func (prof ProfessorDB) Hash() string {
	str := fmt.Sprint(prof.Professor.CodPes)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (prof ProfessorDB) Insert(DB Env) error {
	_, err := DB.Client.Collection("professors").Doc(prof.HashID).Set(DB.Ctx, prof)
	if err != nil {
		return err
	}

	return nil
}
