package domain

import "github.com/google/uuid"

type Task struct {
	Tid    uuid.UUID
	Status string
	Result string
}
