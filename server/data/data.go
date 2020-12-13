package data

/*
	This file is responsible for implementing the REST functions of our most relevant data objects
*/

import (
	"strconv"

	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/utils"
)

// Global variables to store JSON Data
var (
	courses    []entity.Course
	professors []entity.Professor
)

// LoadData is responsible for reading JSON Files and loading to memory
func LoadData() error {

	// I dont know exactly what path makes most sense in our future docker environment, so *rethink later*
	const coursesJSONFileName = "data/courses.json"
	const professorsJSONFileName = "data/professors.json"

	// Read Courses
	var err error
	err = utils.LoadJSON(coursesJSONFileName, &courses)

	if err != nil {
		return err
	}

	// Read Professors
	err = utils.LoadJSON(professorsJSONFileName, &professors)

	if err != nil {
		return err
	}

	return nil
}

// GetProfessors returns list of all professors at every department
func GetProfessors() []entity.Professor {
	return professors
}

// GetProfessorByDepartment returns list of all professors at department 'dep'
func GetProfessorByDepartment(dep string) []entity.Professor {
	return make([]entity.Professor, 0)
}

// GetProfessorByID returns Professor with id 'id'
func GetProfessorByID(id string) entity.Professor {
	for _, professor := range professors {
		if strconv.Itoa(professor.ID) == id {
			return professor
		}
	}
	return entity.Professor{}
}

// GetSubjects returns list of all subjects at every department
func GetSubjects() []entity.Subject {
	return make([]entity.Subject, 0)
}

// GetSubjectByCode returns entity.Subject with code 'code'
func GetSubjectByCode(code string) entity.Subject {
	for _, course := range courses {
		for _, subject := range course.Subjects {
			if subject.Code == code {
				return subject
			}
		}
	}
	return entity.Subject{}
}
