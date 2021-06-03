package private

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

func GetSubjectGrade(ctx *gin.Context, grade *models.Record) {
	ctx.JSON(http.StatusOK, views.NewRecordFromModel(grade))
}

func GetSubjectReview(ctx *gin.Context, review *models.SubjectReview) {
	ctx.JSON(http.StatusOK, views.NewSubjectReviewFromModel(review))
}

func UpdateSubjectReview(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
