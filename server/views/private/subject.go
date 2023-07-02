package private

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

// GetSubjectVerification takes the verification bool and presents it as a response view object
//
// A true verification means that the user has completed the subject
func GetSubjectVerification(ctx *gin.Context, verification bool) {
	ctx.JSON(http.StatusOK, verification)
}

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
