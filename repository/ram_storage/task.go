package ram_storage

import (
	"errors"
	"github.com/google/uuid"
	_ "math/rand"
	_ "strconv"
	"time"
)

type TaskInfo struct {
	Status string
	Result string
}

type Task struct {
	task map[uuid.UUID]TaskInfo
}

func NewTaskInfo() TaskInfo {
	taskInfo := TaskInfo{
		Status: "",
		Result: "",
	}
	return taskInfo
}

func NewTask() *Task {
	return &Task{
		task: make(map[uuid.UUID]TaskInfo),
	}
}

func (rs *Task) Get(uuid uuid.UUID) (*TaskInfo, error) {
	task, exists := rs.task[uuid]
	if !exists {
		return nil, errors.New("task not found")
	}
	return &task, nil
}

func (rs *Task) Put(uuid uuid.UUID, taskStatus string, taskResult string) error {
	task := NewTaskInfo()
	if taskStatus != "" {
		task.Status = taskStatus
	}
	if taskResult != "" {
		task.Result = taskResult
	}
	rs.task[uuid] = task
	return nil
}

func (rs *Task) Post(uuid uuid.UUID, taskStatus string, taskResult string) error {
	if _, exists := rs.task[uuid]; exists {
		return errors.New("task already exists")
	}

	time.Sleep(5 * time.Second)

	task := NewTaskInfo()
	task.Status = taskStatus
	task.Result = taskResult
	rs.task[uuid] = task
	return nil
}

func (rs *Task) Delete(uuid uuid.UUID) error {
	if _, exists := rs.task[uuid]; !exists {
		return errors.New("uuid not found")
	}
	delete(rs.task, uuid)
	return nil
}
