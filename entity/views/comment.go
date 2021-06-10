package views

import (
	"time"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"uuid"`
	Rating    int       `json:"rating"`
	Body      string    `json:"body"`
	Edited    bool      `json:"edited"`
	Timestamp time.Time `json:"timestamp"`
	Upvotes   int       `json:"upvotes"`
	Downvotes int       `json:"downvotes"`
}

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
