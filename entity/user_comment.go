package entity

import (
	"fmt"
	"time"

	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/google/uuid"
)

type UserComment struct {
	ID         uuid.UUID `firestore:"id"`
	Rating     int       `firestore:"rating"`
	Body       string    `firestore:"body"`
	Edited     bool      `firestore:"edited"`
	LastUpdate time.Time `firestore:"last_update"`
	Upvotes    int       `firestore:"upvotes"`
	Downvotes  int       `firestore:"downvotes"`
	Reports    int       `firestore:"reports"`

	Professor      string `firestore:"professor"`
	Subject        string `firestore:"subject"`
	Course         string `firestore:"course"`
	Specialization string `firestore:"specialization"`
}

func (uc UserComment) Hash() string {
	str := fmt.Sprintf("%v%v%v%v", uc.Subject, uc.Course, uc.Specialization, uc.Professor)
	return utils.SHA256(str)
}
