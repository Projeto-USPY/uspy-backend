package views

// CurriculumResult is the response view object for a curriculum query made for a given user
type CurriculumResult struct {
	// Subject data
	Name string `json:"name"`
	Code string `json:"code"`

	// Record data (only present if Completed is true)
	Grade     float64 `json:"grade,omitempty"`
	Frequency int     `json:"frequency,omitempty"`
	Status    string  `json:"status,omitempty"`

	Completed bool `json:"completed"`
}
