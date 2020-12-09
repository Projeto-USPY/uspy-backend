package server

/*
	This file is responsible for implementing the REST functions of our most relevant data objects
*/

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/scraper"
)

// Global variables to store JSON Data
var courses []scraper.Course
var professors []scraper.Professor

// Responsible for reading JSON Files and loading to memory
func LoadData() error {

	// I dont know exactly what path makes most sense in our future docker environment, so *rethink later*
	const coursesJSONFileName = "../data/courses.json"
	const professorsJSONFileName = "../data/professors.json"

	// Read Courses
	coursesJSONFile, err := ioutil.ReadFile(coursesJSONFileName)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(coursesJSONFile), &courses)

	if err != nil {
		return err
	}

	// Read Professors
	professorsJSONFile, err := ioutil.ReadFile(professorsJSONFileName)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(professorsJSONFile), &professors)

	if err != nil {
		return err
	}

	return nil
}

// Todo (return default page)
func DefaultPage(c *gin.Context) {
	c.String(http.StatusOK, "TODO: Default Page")
}

// GetProfessors returns list of all professors at every department
func GetProfessors() []scraper.Professor {
	return professors
}

// GetProfessorByDepartment returns list of all professors at department 'dep'
func GetProfessorByDepartment(dep string) []scraper.Professor {
	return make([]scraper.Professor, 0)
}

// GetProfessorByID returns Professor with id 'id'
func GetProfessorByID(id string) scraper.Professor {
	for _, professor := range professors {
		if strconv.Itoa(professor.ID) == id {
			return professor
		}
	}
	return scraper.Professor{}
}

// GetSubjects returns list of all subjects at every department
func GetSubjects() []scraper.Subject {
	return make([]scraper.Subject, 0)
}

// GetSubjectByCode returns Subject with code 'code'
func GetSubjectByCode(code string) scraper.Subject {
	for _, course := range courses {
		for _, subject := range course.Subjects {
			if subject.Code == code {
				return subject
			}
		}
	}
	return scraper.Subject{}
}
