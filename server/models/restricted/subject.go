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
)

// GetGrades returns all grades from a given subject
func GetGrades(ctx *gin.Context, DB db.Env, sub *controllers.Subject) {
	model := models.NewSubjectFromController(sub)
	snaps, err := DB.RestoreCollection(fmt.Sprintf("subjects/%s/grades", model.Hash()))

	if len(snaps) == 0 {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject grades %v", model))
		return
	} else if err != nil {
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
