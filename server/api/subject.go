package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server/data/course"
	"github.com/tpreischadt/ProjetoJupiter/server/data/subject"
	"net/http"
)

func GetSubjects(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		courses, err := course.GetAll(DB)
		if err != nil {
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

		sub, err := subject.GetByCode(DB, sub.Code)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, sub)
	}
}
