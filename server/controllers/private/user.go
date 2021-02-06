package private

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server/models/private"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func GetSubjectReview(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.MustGet("access_token")
		sub := c.MustGet("Subject").(entity.Subject)

		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		userHash := entity.User{Login: userID}.Hash()
		subHash := entity.Subject{CourseCode: sub.CourseCode, Code: sub.Code}.Hash()

		review, err := private.GetSubjectReview(DB, userHash, subHash)

		if status.Code(err) == codes.NotFound {
			// user has not yet reviewed the subject
			c.Status(http.StatusNotFound)
		} else if err == nil {
			// user has already reviewed the subject
			c.JSON(http.StatusOK, review)
		} else {
			c.Status(http.StatusInternalServerError)
		}
	}
}
