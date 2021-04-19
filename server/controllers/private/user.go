// package private contains the callbacks for every /private endpoint
// for backend-db communication, see /server/models/private
package private

import (
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
	"github.com/Projeto-USPY/uspy-backend/server/models/private"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

// GetSubjectGrade is a closure for the GET /private/subject/grade endpoint
func GetSubjectGrade(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get user and subject info
		token := c.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)
		sub := c.MustGet("Subject").(entity.Subject)

		user, sub := entity.User{Login: userID}, entity.Subject{CourseCode: sub.CourseCode, Code: sub.Code, Specialization: sub.Specialization}
		score, err := private.GetSubjectGrade(DB, user, sub)

		if err == nil {
			c.JSON(http.StatusOK, score)
			return
		}

		c.Status(http.StatusNotFound)
	}
}

// GetSubjectReview is a closure for the GET /private/subject/review endpoint
func GetSubjectReview(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get user and subject info
		token := c.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)
		sub := c.MustGet("Subject").(entity.Subject)

		user, sub := entity.User{Login: userID}, entity.Subject{CourseCode: sub.CourseCode, Code: sub.Code, Specialization: sub.Specialization}

		review, err := private.GetSubjectReview(DB, user, sub)

		if err == nil {
			// user has already reviewed the subject
			c.JSON(http.StatusOK, review)
			return
		}

		// subject does not exist or user has not reviewed it yet
		if err.Error() == "subject does not exist" || status.Code(err) == codes.NotFound {
			c.Status(http.StatusNotFound)
		} else if err.Error() == "user has not done subject" { // user has no permission to review subject
			c.Status(http.StatusForbidden)
		} else {
			log.Println(fmt.Errorf("error fetching review for subject %v, user %v: %v", sub, userID, err))
			c.Status(http.StatusInternalServerError)
		}
	}
}

// UpdateSubjectReview is a closure for the POST /private/subject/review endpoint
func UpdateSubjectReview(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get subject and review data
		sub := c.MustGet("Subject").(entity.Subject)
		sr := entity.SubjectReview{Subject: sub.Code, Course: sub.CourseCode, Specialization: sub.Specialization}

		err := c.ShouldBindJSON(&sr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// get user data
		token := c.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		user := entity.User{Login: userID}

		err = private.UpdateSubjectReview(DB, user, sr)
		if err == nil {
			c.Status(http.StatusOK)
			return
		}

		if err.Error() == "subject does not exist" { // subject doesnt exist
			c.Status(http.StatusNotFound)
		} else if err.Error() == "user has not done subject" { // user has no permission to review subject
			c.Status(http.StatusForbidden)
		} else {
			log.Println(fmt.Errorf("error fetching review for subject %v, user %v: %v", sub, userID, err))
			c.Status(http.StatusInternalServerError)
		}

	}
}
