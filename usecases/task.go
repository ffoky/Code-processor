package usecases

import (
	googleId "github.com/google/uuid"
	"http_server/api/http/types"
	"http_server/repository/ram_storage"
)

type Task interface {
	Get(request types.GetTaskHandlerRequest) (*ram_storage.TaskInfo, error)
	Put(uuid googleId.UUID, status string, result string) error
	Post(uuid googleId.UUID, status string, result string) error
	Delete(uuid googleId.UUID) error
}
