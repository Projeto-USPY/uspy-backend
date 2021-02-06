package public

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

// GetAll returns list of all subjects
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
