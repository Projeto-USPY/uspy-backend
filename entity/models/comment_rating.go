package models

import (
	"github.com/google/uuid"
)

type CommentRating struct {
	ID     uuid.UUID `firestore:"id"`
	Upvote bool      `firestore:"upvote"`
}

func (cr CommentRating) Hash() string {
	return cr.ID.String()
}
