package controllers

type CommentRating struct {
	Offering
	Comment string `form:"comment" binding:"required,uuid"`
}

type CommentVoteBody struct {
	IsUpvote bool `json:"upvote" binding:"required"`
}

type CommentReportBody struct {
	Body string `json:"body" binding:"required,min=10,max=300"`
}
