package controllers

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/tpreischadt/ProjetoJupiter/server/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		sub := c.MustGet("Subject").(entity.Subject)
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
		sub := c.MustGet("Subject").(entity.Subject)

		buckets, err := models.GetGrades(DB, sub)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
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
		ent := c.MustGet("Subject").(entity.Subject)
		sub, err := models.Get(DB, ent)
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

func GetSubjectReview(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.MustGet("access_token")
		sub := c.MustGet("Subject").(entity.Subject)

		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		userHash := entity.User{Login: userID}.Hash()
		subHash := entity.Subject{CourseCode: sub.CourseCode, Code: sub.Code}.Hash()

		// TODO: Check if user has done subject

		snap, err := DB.Restore("users/"+userHash+"/subject_reviews", subHash)
		if status.Code(err) == codes.NotFound {
			// user has not yet reviewed the subject
			c.Status(http.StatusNotFound)
		} else if err == nil {
			// user has already reviewed the subject
			review := entity.SubjectReview{}
			_ = snap.DataTo(&review)
			c.JSON(http.StatusOK, review)
		} else {
			c.Status(http.StatusInternalServerError)
		}

	}
}
