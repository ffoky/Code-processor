package repository

import (
	googleId "github.com/google/uuid"
	"http_server/domain"
)

type Task interface {
	Get(uuid googleId.UUID) (domain.Task, error)
	Put(uuid googleId.UUID, status string, result string) error
	Post(uuid googleId.UUID, status string, result string) error
	Delete(uuid googleId.UUID) error
}
