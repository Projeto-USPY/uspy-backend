package models

import (
	"github.com/google/uuid"
)

type CommentReport struct {
	ID     uuid.UUID `firestore:"id"`
	Report string    `firestore:"report"`

	ProfessorHash  string `firestore:"professor"`
	Subject        string `firestore:"subject"`
	Course         string `firestore:"course"`
	Specialization string `firestore:"specialization"`
}

func (cr CommentReport) Hash() string {
	return cr.ID.String()
}
