package entity

import (
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/middleware"
	"github.com/gin-gonic/gin/binding"
)

// Binders used in controllers for data binding
//
// These are used not only for sanitization purposes, but also for simplifying internal logic
var (
	SubjectBinder         = middleware.Bind("Subject", &controllers.Subject{}, binding.Query)
	OfferingBinder        = middleware.Bind("Offering", &controllers.Offering{}, binding.Query)
	InstituteBinder       = middleware.Bind("Institute", &controllers.Institute{}, binding.Query)
	CourseBinder          = middleware.Bind("Course", &controllers.Course{}, binding.Query)
	CommentRatingBinder   = middleware.Bind("CommentRating", &controllers.CommentRating{}, binding.Query)
	MajorBinder           = middleware.Bind("Major", &controllers.Major{}, binding.Query)
	CurriculumQueryBinder = middleware.Bind("CurriculumQuery", &controllers.CurriculumQuery{}, binding.Query)
	TranscriptQueryBinder = middleware.Bind("TranscriptQuery", &controllers.TranscriptQuery{}, binding.Query)
)
