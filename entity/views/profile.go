package views

// Profile is the response view object for a logged user
//
// It is used to display the Greeting information ("Hello ____")
type Profile struct {
	User string `json:"user"`
	Name string `json:"name"`
}
