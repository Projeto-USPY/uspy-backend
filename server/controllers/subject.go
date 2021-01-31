package controllers

import (
	"errors"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/server/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

func GetSubjects(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		courses, err := models.GetAll(DB)
		if err != nil {
			log.Println(errors.New("error fetching courses: " + err.Error()))
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, courses)
	}
}

func GetSubjectByCode(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		sub := entity.Subject{}
		bindErr := c.BindQuery(&sub)
		if bindErr != nil {
			return
		}

		sub, err := models.Get(DB, sub)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, sub)
	}
}

func GetSubjectGrades(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		sub := entity.Subject{}
		bindErr := c.BindQuery(&sub)
		if bindErr != nil {
			return
		}

		buckets, err := models.GetGrades(DB, sub)
		if err != nil {
			c.Status(http.StatusNotFound)
		}

		avg, approval := 0.0, 0.0
		cnt := 0

		for k, v := range buckets {
			f, _ := strconv.ParseFloat(k, 64)
			avg += f * float64(v)

			if f >= 5.0 {
				approval += float64(v)
			}

			cnt += v
		}

		if len(buckets) == 0 {
			c.Status(http.StatusNotFound)
			return
		}

		avg /= float64(cnt)
		approval /= float64(cnt)

		c.JSON(http.StatusOK, gin.H{"grades": buckets, "average": avg, "approval": approval})
	}
}

func GetSubjectGraph(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		sub := entity.Subject{}
		bindErr := c.BindQuery(&sub)
		if bindErr != nil {
			return
		}

		sub, err := models.Get(DB, sub)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		predecessors, err := models.GetPredecessors(DB, sub)
		if err != nil {
			log.Println(fmt.Errorf("error fetching subject %s/%s predecessors: %s", sub.CourseCode, sub.Code, err.Error()))
			c.Status(http.StatusInternalServerError)
			return
		}

		successors, err := models.GetSuccessors(DB, sub)
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
