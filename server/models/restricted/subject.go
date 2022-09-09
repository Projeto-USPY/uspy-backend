package restricted

import (
	"fmt"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/restricted"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetGrades returns all grades from a given subject
func GetGrades(ctx *gin.Context, DB db.Database, sub *controllers.Subject) {
	model := models.NewSubjectFromController(sub)

	// check subject existence
	if _, err := DB.Restore("subjects/" + model.Hash()); err != nil && status.Code(err) == codes.NotFound {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject %v: %s", model, err.Error()))
		return
	}

	snaps, err := DB.RestoreCollection(fmt.Sprintf("subjects/%s/grades", model.Hash()))

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch subject grades: %s", err.Error()))
		return
	}

	grades := []models.Record{}
	for _, s := range snaps {
		g := models.Record{}
		err := s.DataTo(&g)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind subject %v: %s", model, err.Error()))
			return
		}

		grades = append(grades, g)
	}

	restricted.GetGrades(ctx, grades)
}
