package entity

// Requirement represents a subject requirement
type Requirement struct {
	Subject string `json:"code" firestore:"code"`
	Strong  bool   `json:"strong" firestore:"strong"`
}
