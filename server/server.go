package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/scraper"
)

// Todo (return default page)
func DefaultPage(c *gin.Context) {
	c.String(http.StatusOK, "TODO: Default Page")
}

type Professor struct {
	name string
	dep  string
	id   string
}

// GetProfessors returns list of all professors at every department
func GetProfessors() []Professor {
	return make([]Professor, 0)
}

// GetProfessorByDepartment returns list of all professors at department 'dep'
func GetProfessorByDepartment(dep string) []Professor {
	return make([]Professor, 0)
}

// GetProfessorByID returns Professor with id 'id'
func GetProfessorByID(id string) Professor {
	return Professor{}
}

// GetSubjects returns list of all subjects at every department
func GetSubjects() []scraper.Subject {
	return make([]scraper.Subject, 0)
}

// GetSubjectByCode returns Subject with code 'code'
func GetSubjectByCode(code string) scraper.Subject {
	return scraper.Subject{}
}
