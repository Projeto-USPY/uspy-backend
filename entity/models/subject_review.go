package models

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// SubjectReview is the DTO for a subject review/evaluation made by an user
type SubjectReview struct {
	Subject        string `firestore:"-"`
	Course         string `firestore:"-"`
	Specialization string `firestore:"-"`

	Review map[string]interface{} `firestore:"categories"`
}

// NewSubjectReviewFromController is a constructor. It takes a controller and returns a model.
func NewSubjectReviewFromController(controller *controllers.SubjectReview) *SubjectReview {
	return &SubjectReview{
		Subject:        controller.Subject.Code,
		Course:         controller.Subject.CourseCode,
		Specialization: controller.Subject.Specialization,
		Review:         controller.Review,
	}
}

// Hash returns SHA256(concat(subject, course, specialization))
func (sr SubjectReview) Hash() string {
	str := fmt.Sprintf("%s%s%s", sr.Subject, sr.Course, sr.Specialization)
	return utils.SHA256(str)
}

// Insert sets an offering to a given collection. This is usually /users/#user/subject_reviews
func (sr SubjectReview) Insert(DB db.Database, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(sr.Hash()).Set(DB.Ctx, sr)
	return err
}
