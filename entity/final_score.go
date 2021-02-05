package entity

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
)

type FinalScore struct {
	Grade        float64 `firestore:"grade"`
	Status       string  `firestore:"status"`
	Frequency    int     `firestore:"frequency"`
	OfferingHash string  `firestore:"offering,omitempty"`
}

func (mf FinalScore) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, mf)
	if err != nil {
		return err
	}
	return nil
}
