package storage

import (
	"errors"
	"github.com/google/uuid"
	_ "math/rand"
	_ "strconv"
	"time"
)

type Task struct {
	status string
	result string
}

type RamStorage struct {
	task map[uuid.UUID]Task
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		task: make(map[uuid.UUID]Task),
	}
}

func NewTask() Task {
	task := Task{
		status: "",
		result: "",
	}
	return task
}

func (rs *RamStorage) Get(uuid uuid.UUID) (*Task, error) {
	task, exists := rs.task[uuid]
	if !exists {
		return nil, errors.New("task not found")
	}
	return &task, nil
}

func (rs *RamStorage) Put(uuid uuid.UUID, taskStatus string, taskResult string) error {
	task := NewTask()
	if taskStatus != "" {
		task.status = taskStatus
	}
	if taskResult != "" {
		task.result = taskResult
	}
	rs.task[uuid] = task
	return nil
}

func (rs *RamStorage) Post(uuid uuid.UUID, taskStatus string) error {
	if _, exists := rs.task[uuid]; exists {
		return errors.New("task already exists")
	}

	time.Sleep(15 * time.Second)

	task := NewTask()
	task.status = taskStatus
	rs.task[uuid] = task
	return nil
}

func (rs *RamStorage) Delete(uuid uuid.UUID) error {
	if _, exists := rs.task[uuid]; !exists {
		return errors.New("uuid not found")
	}
	delete(rs.task, uuid)
	return nil
}
