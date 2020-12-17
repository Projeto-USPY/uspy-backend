package db

import (
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

type OfferingDB struct {
	Offering entity.Offering `firestore:"data"`
}

// NewOffering constructor
func NewOffering(ent entity.Offering) *OfferingDB {
	return &OfferingDB{Offering: ent}
}

// md5(concat(subject, professor, year, semester))
func (off OfferingDB) Hash() string {
	concat := fmt.Sprint(
		off.Offering.Subject,
		off.Offering.Professor,
		off.Offering.Year,
		off.Offering.Semester,
	)

	return fmt.Sprintf("%x", md5.Sum([]byte(concat)))
}

func (off OfferingDB) Insert(DB Env) error {
	_, err := DB.Client.Collection("offerings").Doc(off.Hash()).Set(DB.Ctx, off)
	if err != nil {
		return err
	}

	return nil
}
