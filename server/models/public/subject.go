// package public contains functions that implement backend-db communication for every public (not only logged users) /api endpoint
package public

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
	"google.golang.org/api/iterator"
)

// GetByCode returns entity.Subject with code 'code'
// It is the model implementation for /server/models/public.GetSubjectByCode
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

// GetSuccessors returns all subjects that sub is a pre-requisite of.
func GetSuccessors(DB db.Env, sub entity.Subject) (weak, strong []entity.Subject, err error) {
	requirement := entity.Requirement{Subject: sub.Code, Name: sub.Name, Strong: false}
	iter := DB.Client.Collection("subjects").
		Where("true_requirements", "array-contains", requirement).
		Where("course", "==", sub.CourseCode).
		Where("specialization", "==", sub.Specialization).
		Documents(DB.Ctx)

	weak = make([]entity.Subject, 0, 15)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []entity.Subject{}, []entity.Subject{}, err
		}
		var result entity.Subject
		err = doc.DataTo(&result)
		if err != nil {
			return []entity.Subject{}, []entity.Subject{}, err
		}

		weak = append(weak, result)
	}
	iter.Stop()

	requirement = entity.Requirement{Subject: sub.Code, Name: sub.Name, Strong: true}
	iter = DB.Client.Collection("subjects").
		Where("true_requirements", "array-contains", requirement).
		Where("course", "==", sub.CourseCode).
		Where("specialization", "==", sub.Specialization).
		Documents(DB.Ctx)

	strong = make([]entity.Subject, 0, 15)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []entity.Subject{}, []entity.Subject{}, err
		}
		var result entity.Subject
		err = doc.DataTo(&result)
		if err != nil {
			return []entity.Subject{}, []entity.Subject{}, err
		}

		weak = append(strong, result)
	}
	iter.Stop()

	return weak, strong, nil
}
