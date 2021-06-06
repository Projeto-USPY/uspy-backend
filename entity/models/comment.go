package models

import (
	"time"

	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `firestore:"id"`
	Rating    int       `firestore:"rating"`
	Body      string    `firestore:"body"`
	Edited    bool      `firestore:"edited"`
	Timestamp time.Time `firestore:"last_update"`
	Upvotes   int       `firestore:"upvotes"`
	Downvotes int       `firestore:"downvotes"`
	Reports   int       `firestore:"reports"`

	User string `firestore:"-"` // not stored just used for hashing
}

func (c Comment) Hash() string {
	return utils.SHA256(c.User)
}
