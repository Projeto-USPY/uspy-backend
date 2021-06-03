// package controllers
// for backend-db communication, see /server/models/private
package private

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/private"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GetSubjectGrade is a closure for the GET /private/subject/grade endpoint
func GetSubjectGrade(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// get user and subject info
		token := ctx.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)
		sub := ctx.MustGet("Subject").(controllers.Subject)

		private.GetSubjectGrade(ctx, DB, userID, &sub)
	}
}

// GetSubjectReview is a closure for the GET /private/subject/review endpoint
func GetSubjectReview(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// get user and subject info
		token := ctx.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)
		sub := ctx.MustGet("Subject").(controllers.Subject)

		private.GetSubjectReview(ctx, DB, userID, &sub)
	}
}

// UpdateSubjectReview is a closure for the POST /private/subject/review endpoint
func UpdateSubjectReview(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// get subject and review data
		sr := controllers.SubjectReview{}
		err := ctx.ShouldBindJSON(&sr)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// get user data
		token := ctx.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		private.UpdateSubjectReview(ctx, DB, userID, &sr)
	}
}
