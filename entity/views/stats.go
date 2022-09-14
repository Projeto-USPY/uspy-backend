package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

type Stats struct {
	Users     int `json:"users"`
	Grades    int `json:"grades"`
	Subjects  int `json:"subjects"`
	Offerings int `json:"offerings"`
	Comments  int `json:"comments"`
}

// NewStatsFromModel is a constructor. It takes a model stats and returns its response view object.
func NewStatsFromModel(stats *models.Stats) *Stats {
	return &Stats{
		Users:     stats.Users.Count,
		Grades:    stats.Grades.Count,
		Subjects:  stats.Subjects.Count,
		Offerings: stats.Offerings.Count,
		Comments:  stats.Comments.Count,
	}
}
