package utils

import (
	"github.com/Projeto-USPY/uspy-backend/db"
)

func checkSubjectExists(DB db.Env, subHash string) error {
	snap, err := DB.Restore("subjects", subHash)
	if snap == nil || !snap.Exists() {
		return ErrSubjectNotFound
	}
	return err
}

func checkSubjectRecords(DB db.Env, userHash, subHash string) error {
	col, err := DB.RestoreCollection("users/" + userHash + "/final_scores/" + subHash + "/records")
	if len(col) == 0 {
		return ErrNoPermission
	}
	return err
}

// CheckSubjectPermission takes a user hash and a subject hash and checks whether the user has done this subject
func CheckSubjectPermission(DB db.Env, userHash, subHash string) error {
	errSub, errRec := checkSubjectExists(DB, subHash), checkSubjectRecords(DB, userHash, subHash)
	if errSub != nil {
		return errSub
	} else if errRec != nil {
		return errRec
	}

	return nil
}
