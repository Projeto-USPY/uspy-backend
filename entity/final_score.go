package entity

import "github.com/tpreischadt/ProjetoJupiter/db"

type FinalScore struct {
	Grade        int    `firestore:"grade"`
	Status       string `firestore:"status"`
	OfferingHash string `firestore:"offering"`
}

func (mf FinalScore) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, mf)
	if err != nil {
		return err
	}
	return nil
}
