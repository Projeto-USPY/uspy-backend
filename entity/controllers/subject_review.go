package controllers

type SubjectReview struct {
	Subject Subject
	Review  map[string]interface{} `json:"categories" binding:"required,validateSubjectReview"`
}
