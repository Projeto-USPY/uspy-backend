package data

/*
	This file is responsible for implementing the REST functions of our most relevant data objects
*/

import (
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

// GetProfessors returns list of all professors at every department
func GetProfessors() []entity.Professor {
	return make([]entity.Professor, 0)
}

// GetProfessorByDepartment returns list of all professors at department 'dep'
func GetProfessorByDepartment(dep string) []entity.Professor {
	return make([]entity.Professor, 0)
}

// GetProfessorByID returns Professor with id 'id'
func GetProfessorByID(id string) entity.Professor {
	return entity.Professor{}
}

// GetSubjects returns list of all subjects at every department
func GetSubjects() []entity.Subject {
	return make([]entity.Subject, 0)
}

// GetSubjectByCode returns entity.Subject with code 'code'
func GetSubjectByCode(code string) entity.Subject {
	return entity.Subject{}
}
