package usecases

import (
	googleId "github.com/google/uuid"
	"http_server/api/http/types"
	"http_server/domain"
	"net/http"
)

type User interface {
	Get(request types.GetTaskHandlerRequest) (domain.User, error)
	Put(uuid googleId.UUID, login string, password string) error
	Post(login string, password string) error
	Delete(uuid googleId.UUID) error
	Login(login string, password string, w http.ResponseWriter, r *http.Request) (string, error)
}
