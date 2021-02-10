// package public contains the callbacks for every public (not restricted to users) /api endpoint
// for backend-db communication, see /server/models/public
package public

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server/models/public"
	"log"
	"net/http"
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

// GetSubjectByCode is a closure for the GET /api/subject endpoint
func GetSubjectByCode(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		sub := c.MustGet("Subject").(entity.Subject)
		sub, err := public.Get(DB, sub)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, sub)
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

		predecessors, err := public.GetPredecessors(DB, sub)
		if err != nil {
			log.Println(fmt.Errorf("error fetching subject %s/%s predecessors: %s", sub.CourseCode, sub.Code, err.Error()))
			c.Status(http.StatusInternalServerError)
			return
		}

		successors, err := public.GetSuccessors(DB, sub)
		if err != nil {
			log.Println(fmt.Errorf("error fetching subject %s/%s successors: %s", sub.CourseCode, sub.Code, err.Error()))
			c.Status(http.StatusInternalServerError)
			return
		}

		type result struct {
			Code string `json:"code"`
			Name string `json:"name"`
		}
		predecessorsResult := make([]result, 0, 15)

		for i := range predecessors {
			r := result{predecessors[i].Code, predecessors[i].Name}
			predecessorsResult = append(predecessorsResult, r)
		}

		successorsResult := make([]result, 0, 15)

		for i := range successors {
			r := result{successors[i].Code, successors[i].Name}
			successorsResult = append(successorsResult, r)
		}

		c.JSON(http.StatusOK, gin.H{"predecessors": predecessorsResult, "successors": successorsResult})
	}
}