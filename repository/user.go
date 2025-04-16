package repository

import (
	googleId "github.com/google/uuid"
	"http_server/domain"
)

type User interface {
	Get(uuid googleId.UUID) (domain.User, error)
	Put(uuid googleId.UUID, Login string, Password string) error
	Post(uuid googleId.UUID, Login string, Password string) error
	Login(Login string, Password string) (*domain.User, error)
	Delete(uuid googleId.UUID) error
}
