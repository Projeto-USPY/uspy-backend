package public

import (
	"fmt"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/public"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetInstitutes gets all institutes from the database
func GetInstitutes(ctx *gin.Context, DB db.Env) {
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

// GetAllSubjects gets all subjects from the database
func GetAllSubjects(ctx *gin.Context, DB db.Env, controller *controllers.Institute) {
	institute := models.NewInstituteFromController(controller)
	snaps, err := DB.RestoreCollection(fmt.Sprintf(
		"institutes/%s/courses",
		institute.Hash(),
	))

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find collection courses: %s", err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch courses: %s", err.Error()))
		return
	}

	courses := make([]models.Course, 0, 1000)
	for _, s := range snaps {
		var c models.Course
		err = s.DataTo(&c)
		courses = append(courses, c)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch courses: %s", err.Error()))
			return
		}
	}

	public.GetAllSubjects(ctx, courses)
}

// Get gets a subject by its identifier: subject code, course code and course specialization code
func Get(ctx *gin.Context, DB db.Env, sub *controllers.Subject) {
	model := models.Subject{Code: sub.Code, CourseCode: sub.CourseCode, Specialization: sub.Specialization}
	snap, err := DB.Restore("subjects/" + model.Hash())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject %v: %s", model, err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch subject: %s", err))
		return
	}
	err = snap.DataTo(&model)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind subject to object: %s", err))
		return
	}

	public.Get(ctx, &model)
}

// GetRelations gets the subject's graph: their direct predecessors and successors
func GetRelations(ctx *gin.Context, DB db.Env, sub *controllers.Subject) {
	model := models.NewSubjectFromController(sub)
	snap, err := DB.Restore("subjects/" + model.Hash())

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject %v: %s", model, err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch subject: %s", err.Error()))
		return
	}

	if err := snap.DataTo(&model); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not bind subject %v: %s", model, err.Error()))
		return
	}

	getRelatedSubjects := func(model *models.Subject, strength bool) ([]models.Subject, error) {
		requirement := models.Requirement{Subject: model.Code, Name: model.Name, Strong: strength}
		iter := DB.Client.Collection("subjects").
			Where("true_requirements", "array-contains", requirement).
			Where("course", "==", sub.CourseCode).
			Where("specialization", "==", sub.Specialization).
			Documents(DB.Ctx)

		results := make([]models.Subject, 0, 15)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return []models.Subject{}, err
			}
			var result models.Subject
			err = doc.DataTo(&result)
			if err != nil {
				return []models.Subject{}, err
			}

			results = append(results, result)
		}
		iter.Stop()

		return results, nil
	}

	strong, strongErr := getRelatedSubjects(model, true)
	weak, weakErr := getRelatedSubjects(model, false)

	if strongErr != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not get subject predecessors %v: %s", model, strongErr.Error()))
		return
	}

	if weakErr != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not get subject successors %v: %s", model, weakErr.Error()))
		return
	}

	public.GetRelations(ctx, model, weak, strong)
}
