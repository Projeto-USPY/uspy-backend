package private

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

// GetSubjectGrade takes the grade model and presents its response view object
func GetSubjectGrade(ctx *gin.Context, grade *models.Record) {
	ctx.JSON(http.StatusOK, views.NewRecordFromModel(grade))
}

// GetSubjectReview takes the review model and presents its response view object
func GetSubjectReview(ctx *gin.Context, review *models.SubjectReview) {
	ctx.JSON(http.StatusOK, views.NewSubjectReviewFromModel(review))
}

// UpdateSubjectReview is a dummy view method
func UpdateSubjectReview(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
