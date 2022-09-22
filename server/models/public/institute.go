package public

import (
	"fmt"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/public"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetInstitutes gets all institutes from the database
func GetInstitutes(ctx *gin.Context, DB db.Database) {
	snaps, err := DB.RestoreCollection("institutes")
	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find collection institutes: %s", err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch institutes: %s", err.Error()))
		return
	}

	institutes := make([]models.Institute, 0, 200)
	for _, s := range snaps {
		var c models.Institute
		err = s.DataTo(&c)
		institutes = append(institutes, c)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch institutes: %s", err.Error()))
			return
		}
	}

	public.GetInstitutes(ctx, institutes)
}

// GetCourses gets all course codes from a given institute the database
func GetCourses(ctx *gin.Context, DB db.Database, institute *controllers.Institute) {
	model := models.NewInstituteFromController(institute)
	snaps, err := DB.RestoreCollection(fmt.Sprintf(
		"institutes/%s/courses",
		model.Hash(),
	),
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find courses collection from institute: %s", err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch courses from institute: %s", err.Error()))
		return
	}

	courses := make([]*models.Course, 0, 200)
	for _, s := range snaps {
		var c models.Course
		if err := s.DataTo(&c); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch courses: %s", err.Error()))
			return
		}

		c.Subjects = nil // omit subject data to make response payload smaller
		courses = append(courses, &c)
	}

	public.GetCourses(ctx, courses)
}

// GetProfessors gets all course codes from a given institute the database
func GetProfessors(ctx *gin.Context, DB db.Database, institute *controllers.Institute) {
	model := models.NewInstituteFromController(institute)
	snaps, err := DB.RestoreCollection(fmt.Sprintf(
		"institutes/%s/professors",
		model.Hash(),
	),
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find professors collection from institute: %s", err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch professors from institute: %s", err.Error()))
		return
	}

	professors := make([]*models.Professor, 0, 200)
	for _, s := range snaps {
		var c models.Professor
		if err := s.DataTo(&c); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch professors: %s", err.Error()))
			return
		}

		c.CodPesHash = s.Ref.ID
		professors = append(professors, &c)
	}

	public.GetProfessors(ctx, professors)
}
