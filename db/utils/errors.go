package db_utils

import "errors"

var (
	ErrSubjectNotFound = errors.New("subject does not exist")
	ErrNoPermission    = errors.New("user has not done subject")
	ErrCommentNotFound = errors.New("comment not found")
)
