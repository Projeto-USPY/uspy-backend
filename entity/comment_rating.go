package entity

import (
	"github.com/google/uuid"
)

type CommentRating struct {
	ID     uuid.UUID `firestore:"id"`
	Upvote bool      `firestore:"upvote"`
	Report string    `firestore:"report"`
}

func (cr CommentRating) Hash() string {
	return cr.ID.String()
}
