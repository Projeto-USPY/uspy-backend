package db

import (
	"cloud.google.com/go/firestore"
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"golang.org/x/net/context"
)

type OfferingDB struct {
	hashID   string // md5(concat(subject, professor, year, semester)
	offering entity.Offering
}

func calculateHashId(offering entity.Offering) string {
	concat := fmt.Sprintf(
		"%v%v%v%v",
		offering.Subject,
		offering.Professor,
		offering.Year,
		offering.Semester,
	)

	return fmt.Sprintf("%x", md5.Sum([]byte(concat)))
}

// NewOffering constructor
func NewOffering(ent entity.Offering) *OfferingDB {
	return &OfferingDB{hashID: calculateHashId(ent), offering: ent}
}

func (off *OfferingDB) Insert(client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection("offerings").Doc(off.hashID).Set(ctx, off.offering)
	if err != nil {
		return err
	}

	return nil
}
