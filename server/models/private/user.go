package private

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

func GetSubjectReview(DB db.Env, user, sub string) (entity.SubjectReview, error) {
	// TODO: Check if user has done subject
	return entity.SubjectReview{}, nil
}
