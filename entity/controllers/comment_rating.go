package controllers

type CommentRating struct {
	Offering
	ID string `form:"comment" binding:"required,uuid"`
}

type CommentRateBody struct {
	Type string `json:"type" binding:"required,oneof=upvote downvote"`
}

type CommentReportBody struct {
	Body string `json:"body" binding:"required,min=10,max=300"`
}
