package models

import (
	"github.com/google/uuid"
)

type CommentReport struct {
	ID     uuid.UUID `firestore:"id"`
	Report string    `firestore:"report"`
}

func (cr CommentReport) Hash() string {
	return cr.ID.String()
}
