package entity

import "github.com/tpreischadt/ProjetoJupiter/db"

type FinalScore struct {
	Grade        float64 `json:"grade" firestore:"grade"`
	Status       string  `json:"status" firestore:"status"`
	OfferingHash string  `json:"-" firestore:"offering"`
}

func (mf FinalScore) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, mf)
	if err != nil {
		return err
	}
	return nil
}
