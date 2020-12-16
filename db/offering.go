package db

import (
	"cloud.google.com/go/firestore"
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"golang.org/x/net/context"
)

type OfferingDB struct {
	HashID   string          // md5(concat(subject, professor, year, semester)
	Offering entity.Offering `firestore:"data"`
}

// NewOffering constructor
func NewOffering(ent entity.Offering) *OfferingDB {
	off := OfferingDB{Offering: ent}

	off.Hash()
	return &off
}

func (off *OfferingDB) Hash() {
	concat := fmt.Sprint(
		off.Offering.Subject,
		off.Offering.Professor,
		off.Offering.Year,
		off.Offering.Semester,
	)

	off.HashID = fmt.Sprintf("%x", md5.Sum([]byte(concat)))
}

func (off OfferingDB) Insert(client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection("offerings").Doc(off.HashID).Set(ctx, off)
	if err != nil {
		return err
	}

	return nil
}
