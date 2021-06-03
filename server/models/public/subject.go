// package models
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

// GetAll gets all subjects from the database
func GetAll(ctx *gin.Context, DB db.Env) {
	snaps, err := DB.RestoreCollection("courses")
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

	public.GetAll(ctx, courses)
}

// Get gets a subject by its identifier: subject code, course code and course specialization code
func Get(ctx *gin.Context, DB db.Env, sub *controllers.Subject) {
	subModel := models.Subject{Code: sub.Code, CourseCode: sub.CourseCode, Specialization: sub.Specialization}
	snap, err := DB.Restore("subjects", subModel.Hash())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject %v: %s", subModel, err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch subject: %s", err))
		return
	}
	err = snap.DataTo(&sub)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind subject to object: %s", err))
		return
	}

	public.Get(ctx, &subModel)
}

// GetRelations gets the subject's graph: their direct predecessors and successors
func GetRelations(ctx *gin.Context, DB db.Env, sub *controllers.Subject) {
	subModel := models.NewSubjectFromController(sub)
	snap, err := DB.Restore("subjects", subModel.Hash())

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject %v: %s", subModel, err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch subject: %s", err.Error()))
		return
	}

	if err := snap.DataTo(&subModel); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not bind subject %v: %s", subModel, err.Error()))
		return
	}

	getRelatedSubjects := func(subModel *models.Subject, strength bool) ([]models.Subject, error) {
		requirement := models.Requirement{Subject: subModel.Code, Name: subModel.Name, Strong: strength}
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

	strong, strongErr := getRelatedSubjects(subModel, true)
	weak, weakErr := getRelatedSubjects(subModel, false)

	if strongErr != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not get subject predecessors %v: %s", subModel, strongErr.Error()))
		return
	}

	if weakErr != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not get subject successors %v: %s", subModel, weakErr.Error()))
		return
	}

	public.GetRelations(ctx, subModel, weak, strong)
}
