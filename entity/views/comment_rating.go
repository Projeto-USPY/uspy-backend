package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

type CommentRating struct {
	Type string `json:"type"`
}

func NewCommentRatingFromModel(model *models.CommentRating) *CommentRating {
	var responseType string

	if model.Upvote {
		responseType = "upvote"
	} else {
		responseType = "downvote"
	}

	return &CommentRating{Type: responseType}
}
