package domain

import "github.com/google/uuid"

type Task struct {
	ID     uuid.UUID
	Status string
	Result string
}
