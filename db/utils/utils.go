package utils

import (
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
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

// ApplyOperationsInTransaction takes a transaction and a list of operations and applies them sequentially
//
// It returns an error in case any operation fails
// Use this function at the end of a transaction to ensure reads are done in the beginning of the transaction
func ApplyOperationsInTransaction(tx *firestore.Transaction, operators []db.Operation) error {
	for _, obj := range operators {
		if obj.Err != nil {
			return obj.Err
		}

		var operationErr error
		switch obj.Method {
		case "delete":
			operationErr = tx.Delete(obj.Ref)
		case "update":
			operationErr = tx.Update(obj.Ref, obj.Payload.([]firestore.Update))
		case "set":
			operationErr = tx.Set(obj.Ref, obj.Payload)
		default:
			return errors.New(`method not specify, please choose from "set", "delete", "update"`)

		}

		if operationErr != nil {
			return fmt.Errorf("could not apply operation on %#v: %s", obj, operationErr.Error())
		}
	}

	return nil
}
