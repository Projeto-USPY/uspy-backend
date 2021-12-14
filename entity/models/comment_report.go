package models

import (
	"github.com/google/uuid"
)

// CommentReport is the DTO for a comment report, it contains offering data and the report
type CommentReport struct {
	ID     uuid.UUID `firestore:"id"`
	Report string    `firestore:"report"`

	ProfessorHash  string `firestore:"professor"`
	Subject        string `firestore:"subject"`
	Course         string `firestore:"course"`
	Specialization string `firestore:"specialization"`
}

// Hash returns the UUID associated to the reported comment
func (cr CommentReport) Hash() string {
	return cr.ID.String()
}
