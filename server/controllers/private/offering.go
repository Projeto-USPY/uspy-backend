package private

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/private"
	"github.com/gin-gonic/gin"
)

// GetComment is a closure for the GET /private/subjects/offerings/comments endpoint
func GetComment(DB db.Database) func(*gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		off := ctx.MustGet("Offering").(*controllers.Offering)
		off.Subject = *sub

		userID := ctx.MustGet("userID").(string)
		private.GetComment(ctx, DB, userID, off)
	}
}

// GetCommentRating is a closure for the GET /private/subjects/offerings/comments/rating endpoint
func GetCommentRating(DB db.Database) func(*gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		off := ctx.MustGet("Offering").(*controllers.Offering)
		rating := ctx.MustGet("CommentRating").(*controllers.CommentRating)

		off.Subject = *sub
		rating.Offering = *off

		userID := ctx.MustGet("userID").(string)
		private.GetCommentRating(ctx, DB, userID, rating)
	}
}

// RateComment is a closure for the PUT /private/subjects/offerings/comments/rating endpoint
func RateComment(DB db.Database) func(*gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		off := ctx.MustGet("Offering").(*controllers.Offering)
		rating := ctx.MustGet("CommentRating").(*controllers.CommentRating)

		off.Subject = *sub
		rating.Offering = *off

		userID := ctx.MustGet("userID").(string)

		var body controllers.CommentRateBody
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		private.RateComment(ctx, DB, userID, rating, &body)
	}
}

// ReportComment is a closure for the PUT /private/subjects/offerings/comments/report endpoint
func ReportComment(DB db.Database) func(*gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		off := ctx.MustGet("Offering").(*controllers.Offering)
		report := ctx.MustGet("CommentRating").(*controllers.CommentRating)

		off.Subject = *sub
		report.Offering = *off

		userID := ctx.MustGet("userID").(string)

		var body controllers.CommentReportBody
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		private.ReportComment(ctx, DB, userID, report, &body)
	}
}

// PublishComment is a closure for the PUT /private/subjects/offerings/comments endpoint
func PublishComment(DB db.Database) func(*gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		off := ctx.MustGet("Offering").(*controllers.Offering)
		off.Subject = *sub

		userID := ctx.MustGet("userID").(string)

		var comment controllers.Comment
		if err := ctx.ShouldBindJSON(&comment); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		private.PublishComment(ctx, DB, userID, off, &comment)
	}
}

// DeleteComment is a closure for the DELETE /private/subjects/offerings/comments endpoint
func DeleteComment(DB db.Database) func(*gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		off := ctx.MustGet("Offering").(*controllers.Offering)
		off.Subject = *sub

		userID := ctx.MustGet("userID").(string)

		private.DeleteComment(ctx, DB, userID, off)
	}
}
