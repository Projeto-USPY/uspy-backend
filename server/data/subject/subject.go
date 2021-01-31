package subject

import (
	"fmt"

	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"google.golang.org/api/iterator"
)

// GetByCode returns entity.Subject with code 'code'
func Get(DB db.Env, sub entity.Subject) (entity.Subject, error) {
	snap, err := DB.Restore("subjects", sub.Hash())
	if err != nil {
		return entity.Subject{}, err
	}
	err = snap.DataTo(&sub)
	if err != nil {
		return entity.Subject{}, err
	}
	return sub, nil
}

// GetPredecessors returns all subjects which are pre-requisite of sub.
func GetPredecessors(DB db.Env, sub entity.Subject) ([]entity.Subject, error) {

	results := make([]entity.Subject, 0, 15)
	for _, code := range sub.Requirements {
		req := sub
		req.Code = code

		snap, err := DB.Restore("subjects", req.Hash())
		if err != nil {
			return []entity.Subject{}, nil
		}
		err = snap.DataTo(&req)
		if err != nil {
			return []entity.Subject{}, nil
		}

		results = append(results, req)
	}
	return results, nil
}

// GetSuccessors returns all subjects that sub is a pre-requisite of.
func GetSuccessors(DB db.Env, sub entity.Subject) ([]entity.Subject, error) {
	snap, err := DB.Restore("subjects", sub.Hash())
	if err != nil {
		return []entity.Subject{}, nil
	}
	err = snap.DataTo(&sub)
	if err != nil {
		return []entity.Subject{}, nil
	}
	iter := DB.Client.Collection("subjects").
		Where("requirements", "array-contains", sub.Code).
		Where("course", "==", sub.CourseCode).
		Documents(DB.Ctx)

	defer iter.Stop()
	results := make([]entity.Subject, 0, 15)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []entity.Subject{}, err
		}
		var result entity.Subject
		err = doc.DataTo(&result)
		if err != nil {
			return []entity.Subject{}, err
		}

		results = append(results, result)
	}

	return results, nil
}

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
