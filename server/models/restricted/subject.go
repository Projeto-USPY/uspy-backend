// package restricted contains functions that implement backend-db communication for every restricted (logged in users) /api/restricted endpoint
package restricted

import (
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
)

// GetGrades returns all grades from a given subject
func GetGrades(DB db.Env, sub entity.Subject) (map[string]int, error) {
	buckets := make(map[string]int)
	snaps, err := DB.RestoreCollection(fmt.Sprintf("subjects/%s/grades", sub.Hash()))

	if err != nil {
		return map[string]int{}, err
	}

	for _, s := range snaps {
		g := entity.Grade{}
		err := s.DataTo(&g)
		if err != nil {
			return map[string]int{}, err
		}
		buckets[fmt.Sprintf("%.1f", g.Grade)]++
	}

	return buckets, nil
}
