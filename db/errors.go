package db

import "errors"

// DB operation errors
var (
	ErrSubjectNotFound = errors.New("subject does not exist")
	ErrNoPermission    = errors.New("user has not done subject")
	ErrCommentNotFound = errors.New("comment not found")
)
