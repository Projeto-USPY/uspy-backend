package public

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

// GetInstitutes takes the institute models and returns its response view objects
func GetInstitutes(ctx *gin.Context, institutes []models.Institute) {
	viewInstitutes := make([]*views.Institute, 0, len(institutes))
	for _, inst := range institutes {
		viewInstitutes = append(viewInstitutes, views.NewInstituteFromModel(&inst))
	}

	ctx.JSON(http.StatusOK, viewInstitutes)
}

// GetAllSubjects takes the courses models and returns its response view object.
func GetAllSubjects(ctx *gin.Context, courses []models.Course) {
	viewCourses := make([]views.Course, 0, 1000)
	for i := range courses {
		viewCourses = append(viewCourses, *views.NewCourseFromModel(&courses[i]))
	}

	ctx.JSON(http.StatusOK, viewCourses)
}

// Get takes the subject model and returns its response view object
func Get(ctx *gin.Context, model *models.Subject) {
	ctx.JSON(http.StatusOK, views.NewSubjectFromModel(model))
}

// GetRelations takes the subject model and its weak and strong successors
// It returns the view object for its graph
func GetRelations(ctx *gin.Context, model *models.Subject, weak, strong []models.Subject) {
	subView := views.NewSubjectFromModel(model)
	graph := views.SubjectGraph{Predecessors: subView.Requirements}

	for _, w := range weak {
		graph.Successors = append(graph.Successors, views.Requirement{
			Subject: w.Code,
			Name:    w.Name,
			Strong:  false,
		})
	}

	for _, w := range strong {
		graph.Successors = append(graph.Successors, views.Requirement{
			Subject: w.Code,
			Name:    w.Name,
			Strong:  true,
		})
	}

	ctx.JSON(http.StatusOK, graph)
}
