package private

import (
	"errors"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

func GetSubjectReview(DB db.Env, userHash, subHash string) (entity.SubjectReview, error) {
	col, err := DB.RestoreCollection("users/" + userHash + "/final_scores/" + subHash + "/records")
	review := entity.SubjectReview{}

	if len(col) == 0 || err != nil { // user has not done subject
		return review, errors.New("user has not done subject")
	}

	snap, err := DB.Restore("users/"+userHash+"/subject_reviews", subHash)
	if err != nil { // user has not reviewed subject
		return review, err
	}

	err = snap.DataTo(&review)
	return review, err
}
