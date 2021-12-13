package views

import (
	"time"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/google/uuid"
)

// Comment is the response view object for a comment
type Comment struct {
	ID        uuid.UUID `json:"uuid"`
	Rating    int       `json:"rating"`
	Body      string    `json:"body"`
	Edited    bool      `json:"edited"`
	Timestamp time.Time `json:"timestamp"`
	Upvotes   int       `json:"upvotes"`
	Downvotes int       `json:"downvotes"`
}

// NewCommentFromModel is a constructor. It takes a comment model and returns its response view object.
//
// It may panic is the timestamp cannot be generated using the America/Sao_Paulo timezone.
func NewCommentFromModel(model *models.Comment) *Comment {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic("could not time zone America/Sao_Paulo")
	}

	return &Comment{
		ID:        model.ID,
		Rating:    model.Rating,
		Body:      model.Body,
		Edited:    model.Edited,
		Timestamp: model.Timestamp.In(loc),
		Upvotes:   model.Upvotes,
		Downvotes: model.Downvotes,
	}
}
