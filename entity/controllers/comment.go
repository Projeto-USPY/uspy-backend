package controllers

// Comment is the object that holds an offering review comment. It must contain the message and the rating/reaction level.
type Comment struct {
	Rating int    `json:"rating" binding:"required,gte=1,lte=5"`
	Body   string `json:"body" binding:"required,min=10,max=300"`
}
