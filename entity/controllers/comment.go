package controllers

type Comment struct {
	Rating int    `json:"rating" binding:"required,gte=1,lte=5"`
	Body   string `json:"body" binding:"required,min=10,max=300"`
}
