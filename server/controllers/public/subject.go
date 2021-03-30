// package public contains the callbacks for every public (not restricted to users) /api endpoint
// for backend-db communication, see /server/models/public
package public

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
	"github.com/Projeto-USPY/uspy-backend/server/models/public"
	"github.com/gin-gonic/gin"
)

// GetSubject is a closure for the GET /api/subject/all endpoint
func GetSubjects(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		courses, err := public.GetAll(DB)
		if err != nil {
			log.Println(errors.New("error fetching courses: " + err.Error()))
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, courses)
	}
}

// Transforms from map[string][]entity.Requirement to [][]entity.Requirement
func transformRequirements(sub entity.Subject) [][]entity.Requirement {
	requirements := [][]entity.Requirement{}
	for _, val := range sub.Requirements {
		requirements = append(requirements, val)
	}
	return requirements
}

// GetSubjectByCode is a closure for the GET /api/subject endpoint
func GetSubjectByCode(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		sub := c.MustGet("Subject").(entity.Subject)

		sub, err := public.Get(DB, sub)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		type response struct {
			entity.Subject
			Requirements [][]entity.Requirement `json:"requirements"`
		}

		res := response{}
		res.Subject = sub
		res.Requirements = transformRequirements(sub)

		c.JSON(http.StatusOK, res)
	}
}

// GetSubjectGraph is a closure for the GET /api/subject/relations endpoint
func GetSubjectGraph(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		ent := c.MustGet("Subject").(entity.Subject)
		sub, err := public.Get(DB, ent)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		weakSuc, strongSuc, err := public.GetSuccessors(DB, sub)
		if err != nil {
			log.Println(fmt.Errorf("error fetching subject %s/%s successors: %s", sub.CourseCode, sub.Code, err.Error()))
			c.Status(http.StatusInternalServerError)
			return
		}

		type result struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Strong bool   `json:"strong"`
		}

		successorsResult := make([]result, 0, 15)
		for i := range weakSuc {
			r := result{weakSuc[i].Code, weakSuc[i].Name, false}
			successorsResult = append(successorsResult, r)
		}

		for i := range strongSuc {
			r := result{strongSuc[i].Code, strongSuc[i].Name, true}
			successorsResult = append(successorsResult, r)
		}

		c.JSON(http.StatusOK, gin.H{"predecessors": transformRequirements(sub), "successors": successorsResult})
	}
}
