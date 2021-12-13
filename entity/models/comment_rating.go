package models

import (
	"github.com/google/uuid"
)

// CommentRating is the DTO for a comment rating, it contains offering data and the rating evaluation
type CommentRating struct {
	ID     uuid.UUID `firestore:"id"`
	Upvote bool      `firestore:"upvote"`

	ProfessorHash  string `firestore:"professor"`
	Subject        string `firestore:"subject"`
	Course         string `firestore:"course"`
	Specialization string `firestore:"specialization"`
}

// Hash returns the ID associated to the rated comment
func (cr CommentRating) Hash() string {
	return cr.ID.String()
}
