package views

import (
	"time"
)

// Profile is the response view object for a logged user
//
// It is used to display the Greeting information ("Hello ____")
type Profile struct {
	User       string    `json:"user"`
	Name       string    `json:"name"`
	LastUpdate time.Time `json:"last_update"`
}

// NewProfile returns a new view profile object from user data
func NewProfile(user, name string, lastUpdate time.Time) *Profile {
	return &Profile{
		User:       user,
		Name:       name,
		LastUpdate: lastUpdate,
	}
}
