package private

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/private"
	"github.com/gin-gonic/gin"
)

func GetComment(DB db.Env) func(*gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		off := ctx.MustGet("Offering").(*controllers.Offering)
		off.Subject = *sub

		userID := ctx.MustGet("userID").(string)
		private.GetComment(ctx, DB, userID, off)
	}
}

func GetCommentRating(DB db.Env) func(*gin.Context) {
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

func RateComment(DB db.Env) func(*gin.Context) {
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

func ReportComment(DB db.Env) func(*gin.Context) {
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

func PublishComment(DB db.Env) func(*gin.Context) {
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
