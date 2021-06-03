package entity

import (
	"time"

	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/google/uuid"
)

type Comment struct {
	ID         uuid.UUID `firestore:"id" json:"id"`
	Rating     int       `firestore:"rating" json:"rating"`
	Body       string    `firestore:"body" json:"body"`
	Edited     bool      `firestore:"edited" json:"edited"`
	LastUpdate time.Time `firestore:"last_update" json:"timestamp"`
	Upvotes    int       `firestore:"upvotes" json:"upvotes"`
	Downvotes  int       `firestore:"downvotes" json:"downvotes"`
	Reports    int       `firestore:"reports" json:"reports"`

	User string `firestore:"-" json:"-"` // not stored just used for hashing
}

func (c Comment) Hash() string {
	return utils.SHA256(c.User)
}
