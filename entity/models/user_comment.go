package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/utils"
)

// UserComment is the DTO for a user comment
//
// It is a replica of the comment object, but stored in the context of the user instead of the offering
type UserComment struct {
	Comment `firestore:"comment"`

	ProfessorHash  string `firestore:"professor"`
	Subject        string `firestore:"subject"`
	Course         string `firestore:"course"`
	Specialization string `firestore:"specialization"`
}

// Hash returns SHA256(concat(subject, course, specialization, professor_hash))
//
// Note here that professor_hash is a hex sha256 value and not their numeric id.
func (uc UserComment) Hash() string {
	str := fmt.Sprintf("%v%v%v%v", uc.Subject, uc.Course, uc.Specialization, uc.ProfessorHash)
	return utils.SHA256(str)
}
