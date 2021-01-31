package entity

import "github.com/tpreischadt/ProjetoJupiter/db"

type FinalScore struct {
	Grade        float64 `json:"grade" firestore:"grade"`
	Status       string  `json:"status" firestore:"status"`
	Frequency    int     `json:"frequency" firestore:"frequency"`
	OfferingHash string  `json:"-" firestore:"offering,omitempty"`
}

func (mf FinalScore) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, mf)
	if err != nil {
		return err
	}
	return nil
}
