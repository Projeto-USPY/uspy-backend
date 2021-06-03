package models

// Requirement represents a subject requirement
type Requirement struct {
	Subject string `firestore:"code"`
	Name    string `firestore:"name"`
	Strong  bool   `firestore:"strong"`
}
