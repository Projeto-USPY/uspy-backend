package entity

// Subject describes a subject (example: SMA0356 - Cálculo IV)
type Subject struct {
	Code          string
	Name          string
	Description   string
	ClassCredits  int
	AssignCredits int
	TotalHours    string
	Requirements  []string
	Optional      bool
}

// Course represents a course/major (example: BCC)
type Course struct {
	Name     string
	Code     string
	Subjects []Subject
}

// Professor represents a ICMC professor (example: {Moacir Ponti SCC})
type Professor struct {
	ID         int
	Name       string
	Department string
}
