package repository

import (
	googleId "github.com/google/uuid"
	"http_server/repository/ram_storage"
)

type Task interface {
	Get(uuid googleId.UUID) (*ram_storage.TaskInfo, error)
	Put(uuid googleId.UUID, status string, result string) error
	Post(uuid googleId.UUID, status string, result string) error
	Delete(uuid googleId.UUID) error
}
