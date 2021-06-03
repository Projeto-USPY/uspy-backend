package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

type SubjectReview struct {
	Subject        string
	Course         string
	Specialization string

	Review map[string]interface{} `firestore:"categories"`
}

func NewSubjectReviewFromController(controller *controllers.SubjectReview) *SubjectReview {
	return &SubjectReview{
		Subject:        controller.Subject.Code,
		Course:         controller.Subject.CourseCode,
		Specialization: controller.Subject.Specialization,
		Review:         controller.Review,
	}
}

func (sr SubjectReview) Hash() string {
	str := fmt.Sprintf("%s%s%s", sr.Subject, sr.Course, sr.Specialization)
	return utils.SHA256(str)
}

func (sr SubjectReview) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(sr.Hash()).Set(DB.Ctx, sr)
	return err
}
