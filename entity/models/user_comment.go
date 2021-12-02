package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/utils"
)

type UserComment struct {
	Comment `firestore:"comment"`

	ProfessorCode  string `firestore:"professor"`
	Subject        string `firestore:"subject"`
	Course         string `firestore:"course"`
	Specialization string `firestore:"specialization"`
}

func (uc UserComment) Hash() string {
	str := fmt.Sprintf("%v%v%v%v", uc.Subject, uc.Course, uc.Specialization, uc.ProfessorCode)
	return utils.SHA256(str)
}
