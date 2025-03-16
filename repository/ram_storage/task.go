package ram_storage

import (
	"github.com/google/uuid"
	"http_server/domain"
	"http_server/repository"
)

type Task struct {
	tasks map[uuid.UUID]domain.Task
}

func NewTask() *Task {
	return &Task{
		tasks: make(map[uuid.UUID]domain.Task),
	}
}

func (rs *Task) Get(id uuid.UUID) (domain.Task, error) {
	task, exists := rs.tasks[id]
	if !exists {
		return domain.Task{}, repository.NotFound
	}
	return task, nil
}

func (rs *Task) Put(id uuid.UUID, status string, result string) error {
	task, exists := rs.tasks[id]
	if !exists {
		return repository.NotFound
	}

	if status != "" {
		task.Status = status
	}
	if result != "" {
		task.Result = result
	}

	rs.tasks[id] = task
	return nil
}

func (rs *Task) Post(id uuid.UUID, status string, result string) error {
	if _, exists := rs.tasks[id]; exists {
		return repository.ErrTaskExists
	}

	rs.tasks[id] = domain.Task{
		ID:     id,
		Status: status,
		Result: result,
	}
	return nil
}

func (rs *Task) Delete(id uuid.UUID) error {
	if _, exists := rs.tasks[id]; !exists {
		return repository.NotFound
	}
	delete(rs.tasks, id)
	return nil
}
