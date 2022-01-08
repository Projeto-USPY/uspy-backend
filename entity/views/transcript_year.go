package views

// TranscriptYear is the response view object for a when the user's transcript years is queried
type TranscriptYear struct {
	Year      int   `json:"year"`
	Semesters []int `json:"semesters"`
}
