package public

import (
	"fmt"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/server/views/public"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAllSubjects gets all subjects from a given course in the database
func GetAllSubjects(ctx *gin.Context, DB db.Database, controller *controllers.Course) {
	course := models.NewCourseFromController(controller)
	snap, err := DB.Restore(fmt.Sprintf(
		"institutes/%s/courses/%s",
		utils.SHA256(course.Institute),
		course.Hash(),
	))

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find course: %s", err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch course: %s", err.Error()))
		return
	}

	var courseModel models.Course
	if err := snap.DataTo(&courseModel); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind course to object: %s", err))
		return
	}

	public.GetAllSubjects(ctx, &courseModel)
}

// Get gets a subject by its identifier: subject code, course code and course specialization code
func Get(ctx *gin.Context, DB db.Database, sub *controllers.Subject) {
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
func GetRelations(ctx *gin.Context, DB db.Database, sub *controllers.Subject) {
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

// GetSiblingSubjects gets the subject's siblings: Subjects with same institute, course, specialization and semester
func GetSiblingSubjects(ctx *gin.Context, DB db.Database, sub *controllers.Subject) {
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

	iter := DB.Client.Collection("subjects").Query.
		Where("course", "==", model.CourseCode).
		Where("specialization", "==", model.Specialization).
		Where("semester", "==", model.Semester).
		Documents(DB.Ctx)

	siblings := make([]*views.SubjectSibling, 0, 15)
	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch subject: %s", err.Error()))
			return
		}

		var subject models.Subject
		if err := snap.DataTo(&subject); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind subject to object: %s", err.Error()))
			return
		}

		siblings = append(siblings, views.NewSubjectSibling(&subject))
	}

	public.GetSiblingSubjects(ctx, siblings)
}
