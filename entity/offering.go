package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
)

// Offering describes an offering of a subject (example: Cálculo IV - 2019.2)
type Offering struct {
	Semester  int    `firestore:"semester,omitempty"`
	Year      int    `firestore:"year,omitempty"`
	Professor int    `firestore:"professor,omitempty"`
	Subject   string `firestore:"subject,omitempty"`
}

// md5(concat(subject, professor, year, semester))
func (off Offering) Hash() string {
	concat := fmt.Sprint(
		off.Subject,
		off.Professor,
		off.Year,
		off.Semester,
	)

	return fmt.Sprintf("%x", md5.Sum([]byte(concat)))
}

func (off Offering) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(off.Hash()).Set(DB.Ctx, off)
	if err != nil {
		return err
	}

	return nil
}
