package views

type SubjectGraph struct {
	Predecessors [][]Requirement `json:"predecessors"`
	Successors   []Requirement   `json:"successors"`
}
