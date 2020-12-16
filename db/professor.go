package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"crypto/md5"
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
func NewProfessor(ent entity.Professor, client *firestore.Client, ctx context.Context) (*ProfessorDB, error) {
	col := client.Collection("offerings")
	offs := col.Where("data.professor", "==", ent.CodPes)
	iter := offs.Documents(ctx)
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

func (prof *ProfessorDB) Hash() {
	str := fmt.Sprint(prof.Professor.CodPes)
	prof.HashID = fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (prof *ProfessorDB) Insert(client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection("professors").Doc(prof.HashID).Set(ctx, prof)
	if err != nil {
		return err
	}

	return nil
}
