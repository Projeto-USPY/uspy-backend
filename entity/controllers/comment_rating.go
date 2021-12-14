package controllers

// CommentRating is the object used for identifying the comment being rated
type CommentRating struct {
	Offering
	ID string `form:"comment" binding:"required,uuid"`
}

// CommentRateBody is the object that holds the type when rating other users' comments
type CommentRateBody struct {
	Type string `json:"type" binding:"required,oneof=upvote downvote none"`
}

// CommentReportBody is the object that holds the message when reporting users' comments
type CommentReportBody struct {
	Body string `json:"body" binding:"required,min=10,max=300"`
}
