// package public contains functions that implement backend-db communication for every public (not only logged users) /api endpoint
package public

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
)

// GetAll returns list of all subjects
// GetAll is the model implementation for /server/controllers/public.GetSubjects
func GetAll(DB db.Env) ([]entity.Course, error) {
	snaps, err := DB.RestoreCollection("courses")
	if err != nil {
		return []entity.Course{}, err
	}
	courses := make([]entity.Course, 0, 1000)
	for _, s := range snaps {
		c := entity.Course{}
		err = s.DataTo(&c)
		courses = append(courses, c)
		if err != nil {
			return []entity.Course{}, err
		}
	}
	return courses, nil
}
