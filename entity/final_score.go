package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
)

type FinalScore struct {
	Grade        float64 `firestore:"grade"`
	Status       string  `firestore:"status"`
	Year         int     `firestore:"-"`
	Semester     int     `firestore:"-"`
	Frequency    int     `firestore:"frequency"`
	OfferingHash string  `firestore:"offering,omitempty"`
}

func (mf FinalScore) Hash() string {
	str := fmt.Sprintf("%d%d", mf.Year, mf.Semester)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (mf FinalScore) Insert(DB db.Env, collection string) error {
	_, _, err := DB.Client.Collection(collection).Add(DB.Ctx, mf)
	if err != nil {
		return err
	}
	return nil
}
