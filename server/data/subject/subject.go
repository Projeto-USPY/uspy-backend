package subject

import (
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

// GetByCode returns entity.Subject with code 'code'
func GetByCode(DB db.Env, code string) (entity.Subject, error) {
	subject := entity.Subject{Code: code}
	snap, err := DB.Restore("subjects", subject.Hash())
	if err != nil {
		return entity.Subject{}, err
	}
	err = snap.DataTo(&subject)
	if err != nil {
		return entity.Subject{}, err
	}
	return subject, nil
}

// GetGrades returns all grades from a given subject
func GetGrades(DB db.Env, code string) (map[string]int, error) {
	buckets := make(map[string]int)
	subject := entity.Subject{Code: code}
	snaps, err := DB.RestoreCollection(fmt.Sprintf("subjects/%s/grades", subject.Hash()))
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
