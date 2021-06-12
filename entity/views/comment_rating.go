package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

type CommentRating struct {
	IsUpvote bool `json:"upvote"`
}

func NewCommentRatingFromModel(model *models.CommentRating) *CommentRating {
	return &CommentRating{IsUpvote: model.Upvote}
}
