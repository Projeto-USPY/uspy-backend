package private

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

func GetSubjectGrade(ctx *gin.Context, grade *models.FinalScore) {
	ctx.JSON(http.StatusOK, views.NewFinalScoreFromModel(grade))
}

func GetSubjectReview(ctx *gin.Context, review *models.SubjectReview) {
	ctx.JSON(http.StatusOK, views.NewSubjectReviewFromModel(review))
}

func UpdateSubjectReview(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
