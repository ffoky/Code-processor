package usecases

import (
	googleId "github.com/google/uuid"
	"http_server/api/http/types"
	"http_server/domain"
)

type Task interface {
	Get(request types.GetTaskHandlerRequest) (domain.Task, error)
	Put(uuid googleId.UUID, status string, result string) error
	Post() (googleId.UUID, error)
	Delete(uuid googleId.UUID) error
}
