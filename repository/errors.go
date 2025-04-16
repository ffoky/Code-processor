package repository

import "errors"

var (
	NotFound         = errors.New("uuid not found")
	InternalError    = errors.New("failed to complete task")
	ErrTaskExists    = errors.New("task already exists")
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("login or password is incorrect")
	ErrSessionExists = errors.New("session already exists")
)
