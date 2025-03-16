package service

import (
	googleId "github.com/google/uuid"
	"http_server/api/http/types"
	"http_server/domain"
	"http_server/repository"
	"math/rand"
	_ "net/http"
	"strconv"
	"time"
)

type Task struct {
	repo repository.Task
}

func NewTask(repo repository.Task) *Task {
	return &Task{
		repo: repo,
	}
}

func (rs *Task) Get(req types.GetTaskHandlerRequest) (domain.Task, error) {
	return rs.repo.Get(req.Uuid)
}

func (rs *Task) Put(uuid googleId.UUID, status string, result string) error {
	return rs.repo.Put(uuid, status, result)
}

func (rs *Task) generateUUID() googleId.UUID {
	id := googleId.New()
	return id
}

func (rs *Task) CompleteTask(uuid googleId.UUID) error {
	time.Sleep(25 * time.Second)
	taskStatus := "ready"
	taskResult := strconv.Itoa(rand.Intn(100))
	if err := rs.Put(uuid, taskStatus, taskResult); err != nil {
		return repository.InternalError
	}
	return nil
}

func (rs *Task) Post() (googleId.UUID, error) {
	uuid := rs.generateUUID()
	status := "in_progress"
	result := ""
	go func() {
		err := rs.CompleteTask(uuid)
		if err != nil {

		}
	}()
	return uuid, rs.repo.Post(uuid, status, result)
}

func (rs *Task) Delete(uuid googleId.UUID) error {
	return rs.repo.Delete(uuid)
}
