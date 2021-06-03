package views

// Requirement represents a subject requirement
type Requirement struct {
	Subject string `json:"code"`
	Name    string `json:"name"`
	Strong  bool   `json:"strong"`
}
