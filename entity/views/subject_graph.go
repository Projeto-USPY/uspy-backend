package views

// SubjectGraph is the response view object for a subjects predecessors and successors graph
//
// Predecessors represents the subjects that need to be taken before a given subject
// Successors represents the successors that need to be taken after a given subject
type SubjectGraph struct {
	Predecessors [][]Requirement `json:"predecessors"`
	Successors   []Requirement   `json:"successors"`
}
