package utils

import (
	"errors"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
)

func checkSubjectExists(DB db.Env, subHash string) error {
	snap, err := DB.Restore("subjects/" + subHash)
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

// ApplyConcurrentOperationsInTransaction takes a transaction and a list of operations and applies them using concurrency
//
// It returns an error in case any operation fails
// Use this function at the end of a transaction to ensure write operations are done in the end of the transaction
func ApplyConcurrentOperationsInTransaction(tx *firestore.Transaction, operators []db.Operation) error {
	var wg sync.WaitGroup

	// launch producers
	errChan := make(chan error, len(operators))
	wg.Add(len(operators))
	for _, obj := range operators {
		go func(job db.Operation, wg *sync.WaitGroup) {
			defer wg.Done()

			if job.Err != nil { // if any operator failed, stop all goroutines, then return error
				errChan <- job.Err
				return
			}

			var operationErr error
			switch job.Method {
			case "delete":
				operationErr = tx.Delete(job.Ref)
			case "update":
				operationErr = tx.Update(job.Ref, job.Payload.([]firestore.Update))
			case "set":
				operationErr = tx.Set(job.Ref, job.Payload)
			default: // incorrect operation
				operationErr = errors.New(`method not specified, please choose from "set", "delete", "update"`)
			}

			errChan <- operationErr
		}(obj, &wg)
	}

	wg.Wait()
	close(errChan)

	// receiver statuses
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
