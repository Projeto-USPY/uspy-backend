package controllers

type SubjectReview struct {
	Subject
	Review map[string]interface{} `json:"categories" binding:"required,validateSubjectReview"`
}
