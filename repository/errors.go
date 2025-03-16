package repository

import "errors"

var (
	NotFound      = errors.New("uuid not found")
	InternalError = errors.New("failed to complete task")
	ErrTaskExists = errors.New("task already exists")
)
