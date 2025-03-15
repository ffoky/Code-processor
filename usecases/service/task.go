package service

import (
	googleId "github.com/google/uuid"
	"http_server/api/http/types"
	"http_server/repository"
	"http_server/repository/ram_storage"
)

type Task struct {
	repo repository.Task
}

func NewTask(repo repository.Task) *Task {
	return &Task{
		repo: repo,
	}
}

func (rs *Task) Get(req types.GetTaskHandlerRequest) (*ram_storage.TaskInfo, error) {
	return rs.repo.Get(req.Uuid)
}

func (rs *Task) Put(uuid googleId.UUID, status string, result string) error {
	return rs.repo.Put(uuid, status, result)
}

func (rs *Task) Post(uuid googleId.UUID, status string, result string) error {
	return rs.repo.Post(uuid, status, result)
}

func (rs *Task) Delete(uuid googleId.UUID) error {
	return rs.repo.Delete(uuid)
}
