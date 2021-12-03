package models

import (
	"github.com/google/uuid"
)

type CommentRating struct {
	ID     uuid.UUID `firestore:"id"`
	Upvote bool      `firestore:"upvote"`

	ProfessorHash  string `firestore:"professor"`
	Subject        string `firestore:"subject"`
	Course         string `firestore:"course"`
	Specialization string `firestore:"specialization"`
}

func (cr CommentRating) Hash() string {
	return cr.ID.String()
}
