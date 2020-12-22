package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server/data/course"
	"github.com/tpreischadt/ProjetoJupiter/server/data/subject"
)

func GetSubjects(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		courses, err := course.GetAll(DB)
		if err != nil {
			log.Print(err.Error())
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

		sub, err := subject.Get(DB, sub)
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

		buckets, err := subject.GetGrades(DB, sub)
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

		c.JSON(http.StatusOK, gin.H{"Grades": buckets, "Average": avg, "Approval": approval})
	}
}

func GetSubjectGraph(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		sub := entity.Subject{}
		bindErr := c.BindQuery(&sub)
		if bindErr != nil {
			return
		}

		sub, err := subject.Get(DB, sub)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		successors, err := subject.GetSucessors(DB, sub)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		type result struct {
			Code string
			Name string
		}
		results := make([]result, 0, 20)

		for i := range successors {
			r := result{successors[i].Code, successors[i].Name}
			results = append(results, r)
		}

		c.JSON(http.StatusOK, gin.H{"Predecessors": sub.Requirements, "Successors": results})
	}
}
