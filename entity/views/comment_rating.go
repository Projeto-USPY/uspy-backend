package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

// CommentRating is the response view object for a comment rating.
type CommentRating struct {
	Type string `json:"type"`
}

// NewCommentRatingFromModel is a constructor. It takes a comment rating model and returns its response view object.
func NewCommentRatingFromModel(model *models.CommentRating) *CommentRating {
	var responseType string

	if model.Upvote {
		responseType = "upvote"
	} else {
		responseType = "downvote"
	}

	return &CommentRating{Type: responseType}
}
