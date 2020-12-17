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

// Offering describes an offering of a subject (example: Cálculo IV - 2019.2)
type Offering struct {
	Semester  int    `firestore:"semester"`
	Year      int    `firestore:"year"`
	Professor int    `firestore:"professor"`
	Subject   string `firestore:"subject"`
}

// TODO: Change subject/course entity or add another collection to DB?
// Course represents a course/major (example: BCC)
type Course struct {
	Name     string
	Code     string
	Subjects []Subject
}

// Professor represents a ICMC professor (example: {Moacir Ponti SCC})
type Professor struct {
	CodPes     int    `firestore:"code"`
	Name       string `firestore:"name"`
	Department string `firestore:"dep"`
}

// User represents an user
type User struct {
	ID       int
	Login    string `json:"login" binding:"required"`
	Password string `json:"pwd" binding:"required"` // used only because of REST requests, do not store in db
}
