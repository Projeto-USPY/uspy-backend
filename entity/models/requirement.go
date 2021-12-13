package models

// Requirement represents a subject requirement
//
// It is nested inside a subject DTO and not mapped directly to a document
type Requirement struct {
	Subject string `firestore:"code"`
	Name    string `firestore:"name"`
	Strong  bool   `firestore:"strong"`
}
