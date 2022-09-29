package models

// ProfessorReview is the DTO for a professor review/evaluation made by an user
type ProfessorReview struct {
	Review map[string]bool `firestore:"categories"`
}
