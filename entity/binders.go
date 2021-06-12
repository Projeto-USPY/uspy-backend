package entity

import (
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/middleware"
	"github.com/gin-gonic/gin/binding"
)

var (
	SubjectBinder       = middleware.Bind("Subject", &controllers.Subject{}, binding.Query)
	OfferingBinder      = middleware.Bind("Offering", &controllers.Offering{}, binding.Query)
	CommentRatingBinder = middleware.Bind("CommentRating", &controllers.CommentRating{}, binding.Query)
)
