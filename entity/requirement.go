package entity

// Requirement represents a subject requirement
type Requirement struct {
	Subject string `json:"code" firestore:"code"`
	Name    string `json:"name" firestore:"name"`
	Strong  bool   `json:"strong" firestore:"strong"`
}
