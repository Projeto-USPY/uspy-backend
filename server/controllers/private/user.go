package private

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server/models/private"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

func GetSubjectReview(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.MustGet("access_token")
		sub := c.MustGet("Subject").(entity.Subject)

		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		user, sub := entity.User{Login: userID}, entity.Subject{CourseCode: sub.CourseCode, Code: sub.Code}

		review, err := private.GetSubjectReview(DB, user, sub)

		if err == nil {
			// user has already reviewed the subject
			c.JSON(http.StatusOK, review)
			return
		}

		if status.Code(err) == codes.NotFound {
			// user has not yet reviewed the subject or the subject doesnt exist
			c.Status(http.StatusNotFound)
		} else if err.Error() == "user has not done subject" {
			c.Status(http.StatusForbidden)
		} else {
			log.Println(fmt.Errorf("error fetching review for subject %v, user %v: %v", sub, userID, err))
			c.Status(http.StatusInternalServerError)
		}
	}
}

func UpdateSubjectReview(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		sub := c.MustGet("Subject").(entity.Subject)
		sr := entity.SubjectReview{Subject: sub.Code, Course: sub.CourseCode}

		err := c.ShouldBindJSON(&sr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		token := c.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		user := entity.User{Login: userID}

		err = private.UpdateSubjectReview(DB, user, sr)
		if err != nil {
			log.Println(fmt.Errorf("error updating subject review: " + err.Error()))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}
