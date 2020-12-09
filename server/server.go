package server

/*
	This file is responsible for implementing the REST functions of our most relevant data objects
*/

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/scraper"
)

// Todo (return default page)
func DefaultPage(c *gin.Context) {
	c.String(http.StatusOK, "TODO: Default Page")
}

// GetProfessors returns list of all professors at every department
func GetProfessors() []scraper.Professor {
	return make([]scraper.Professor, 0)
}

// GetProfessorByDepartment returns list of all professors at department 'dep'
func GetProfessorByDepartment(dep string) []scraper.Professor {
	return make([]scraper.Professor, 0)
}

// GetProfessorByID returns Professor with id 'id'
func GetProfessorByID(id string) scraper.Professor {
	return scraper.Professor{}
}

// GetSubjects returns list of all subjects at every department
func GetSubjects() []scraper.Subject {
	return make([]scraper.Subject, 0)
}

// GetSubjectByCode returns Subject with code 'code'
func GetSubjectByCode(code string) scraper.Subject {
	return scraper.Subject{}
}
