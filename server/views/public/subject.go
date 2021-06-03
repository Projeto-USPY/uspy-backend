package public

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

func GetAll(ctx *gin.Context, courses []models.Course) {
	viewCourses := make([]views.Course, 0, 1000)
	for i := range courses {
		viewCourses = append(viewCourses, *views.NewCourseFromModel(&courses[i]))
	}

	ctx.JSON(http.StatusOK, viewCourses)
}

func Get(ctx *gin.Context, subModel *models.Subject) {
	ctx.JSON(http.StatusOK, views.NewSubjectFromModel(subModel))
}

func GetRelations(ctx *gin.Context, subModel *models.Subject, weak, strong []models.Subject) {
	subView := views.NewSubjectFromModel(subModel)
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
