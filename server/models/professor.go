package models

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

// GetProfessors returns list of all professors at every department
func GetProfessors(DB db.Env) ([]entity.Professor, error) {
	snaps, err := DB.Client.Collection("professors").Documents(DB.Ctx).GetAll()
	if err != nil {
		return []entity.Professor{}, err
	}
	profs := make([]entity.Professor, 0, 200)
	for _, d := range snaps {
		prof := entity.Professor{}
		err := d.DataTo(&prof)
		if err != nil {
			return []entity.Professor{}, err
		}
		profs = append(profs, prof)
	}

	return profs, nil
}

// GetProfessorByDepartment returns list of all professors at department 'dep'
func GetProfessorByDepartment(dep string) []entity.Professor {
	return make([]entity.Professor, 0)
}

// GetProfessorByID returns Professor with id 'id'
func GetProfessorByID(id string) entity.Professor {
	return entity.Professor{}
}
