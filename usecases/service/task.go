package service

import (
	"fmt"
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
	repo   repository.Task
	sender repository.ObjectSender
}

func NewTask(repo repository.Task, sender repository.ObjectSender) *Task {
	return &Task{
		repo:   repo,
		sender: sender,
	}
}

func (rs *Task) Get(req types.GetTaskHandlerRequest) (domain.Task, error) {
	return rs.repo.Get(req.Uuid)
}

func (rs *Task) Put(task domain.Task) error {
	return rs.repo.Put(task.Tid, task.Status, task.Result)
}

func (rs *Task) generateUUID() googleId.UUID {
	id := googleId.New()
	return id
}

func (rs *Task) CompleteTask(uuid googleId.UUID) error {
	time.Sleep(25 * time.Second)
	taskStatus := "ready"
	taskResult := strconv.Itoa(rand.Intn(100))
	task := domain.Task{Tid: uuid, Status: taskStatus, Result: taskResult}
	if err := rs.Put(task); err != nil {
		return repository.InternalError
	}
	return nil
}

func (rs *Task) Post() (googleId.UUID, error) { //TODO переписать аргумент как domain.Task
	uuid := rs.generateUUID()
	status := "in_progress"
	result := ""
	go func() {
		err := rs.CompleteTask(uuid)
		if err != nil {
			//TODO обработать ошибку
		}
	}()
	err := rs.sender.Send(domain.Task{Tid: uuid, Status: status, Result: result})
	if err != nil {
		return googleId.UUID{}, fmt.Errorf("sending object: %w", err)
	}
	return uuid, rs.repo.Post(uuid, status, result)
}

func (rs *Task) Delete(uuid googleId.UUID) error {
	return rs.repo.Delete(uuid)
}
