package models

import (
	"time"

	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/google/uuid"
)

// Comment is the DTO for a comment
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

// Hash returns SHA256(userID), where userID is the the ID of the user who made the comment
//
// This is made this way so that looking up comments is easy and fast.
func (c Comment) Hash() string {
	return utils.SHA256(c.User)
}
