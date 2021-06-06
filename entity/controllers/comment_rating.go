package controllers

type CommentRating struct {
	Offering Offering
	Comment  string `form:"comment" binding:"required,len=32,alphanum"`
	IsUpvote bool   `json:"upvote" binding:"required"`
}
