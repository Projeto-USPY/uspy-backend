package views

import (
	"time"

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
