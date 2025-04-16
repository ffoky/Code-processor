package domain

import "github.com/google/uuid"

type User struct {
	Uid      uuid.UUID
	Login    string
	Password string
}
