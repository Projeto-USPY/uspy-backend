// package models
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
func GetGrades(ctx *gin.Context, DB db.Env, sub *controllers.Subject) {
	subModel := models.NewSubjectFromController(sub)
	snaps, err := DB.RestoreCollection(fmt.Sprintf("subjects/%s/grades", subModel.Hash()))

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject %v: %s", subModel, err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch subject: %s", err.Error()))
		return
	}

	grades := []models.Grade{}
	for _, s := range snaps {
		g := models.Grade{}
		err := s.DataTo(&g)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind subject %v: %s", subModel, err.Error()))
			return
		}

		grades = append(grades, g)
	}

	restricted.GetGrades(ctx, grades)
}
