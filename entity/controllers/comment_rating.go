package controllers

type CommentRating struct {
	Offering Offering
	Comment  string `form:"comment" binding:"required,uuid"`
	IsUpvote bool   `json:"upvote" binding:"required"`
}
