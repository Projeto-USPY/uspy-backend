package models

import (
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// StatsEntry represents a simple counter
type StatsEntry struct {
	Name  string `firestore:"name"`
	Count int    `firestore:"count"`

	LastUpdate time.Time `firestore:"last_update"`
}

// Stats is a DTO for the database stats
type Stats struct {
	Users     StatsEntry `firestore:"users"`
	Grades    StatsEntry `firestore:"grades"`
	Subjects  StatsEntry `firestore:"subjects"`
	Offerings StatsEntry `firestore:"offerings"`
	Comments  StatsEntry `firestore:"comments"`
}

func (s Stats) Hash() string {
	return utils.SHA256("uspy")
}

// Update updates the user and grades count
//
// Other data is not updated here to avoid heavy operations
// Use uspy-scraper to sync count of scraped data
func (s Stats) Update(DB db.Database, collection string) error {
	updates := make([]firestore.Update, 0)

	updates = append(updates,
		firestore.Update{
			Path:  "users.count",
			Value: s.Users.Count,
		},
		firestore.Update{
			Path:  "users.last_update",
			Value: s.Users.LastUpdate,
		},
		firestore.Update{
			Path:  "grades.count",
			Value: s.Grades.Count,
		},
		firestore.Update{
			Path:  "grades.last_update",
			Value: s.Grades.LastUpdate,
		},
	)

	_, err := DB.Client.Collection(collection).Doc(s.Hash()).Update(DB.Ctx, updates)
	return err
}

// Insert sets the stats counters
func (s Stats) Insert(DB db.Database, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(s.Hash()).Set(DB.Ctx, s)
	return err
}
