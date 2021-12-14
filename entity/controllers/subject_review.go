package controllers

// SubjectReview is the object used for evaluating a subject, without the context of a professor
type SubjectReview struct {
	Subject
	Review map[string]interface{} `json:"categories" binding:"required,validateSubjectReview"`
}
