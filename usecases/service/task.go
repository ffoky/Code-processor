package service

import (
	"fmt"
	googleId "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"http_server/api/http/types"
	"http_server/domain"
	"http_server/repository"
	"math/rand"
	_ "net/http"
	"strconv"
	"sync"
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
		logrus.Errorf("Service failed to comlete task %v", err)
		return repository.InternalError
	}
	return nil
}

func (rs *Task) Post() (googleId.UUID, error) {
	uuid := rs.generateUUID()
	status := "in_progress"
	result := ""
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := rs.CompleteTask(uuid)
		if err != nil {
			logrus.Errorf("failed completing task: %v", err)
		}
	}()
	err := rs.sender.Send(domain.Task{Tid: uuid, Status: status, Result: result})
	if err != nil {
		logrus.Errorf("sending object: %v", err)
		return googleId.UUID{}, fmt.Errorf("sending object: %w", err)
	}
	return uuid, rs.repo.Post(uuid, status, result)
}

func (rs *Task) Delete(uuid googleId.UUID) error {
	return rs.repo.Delete(uuid)
}
