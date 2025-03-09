package storage

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

type RamStorage struct {
	task map[string]string
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		task: make(map[string]string),
	}
}

func (rs *RamStorage) Get(uuid string) (*string, error) {
	taskStatusResult, exists := rs.task[uuid]
	if !exists {
		return nil, nil
	}
	return &taskStatusResult, nil
}

func (rs *RamStorage) Put(uuid string, status string) error {
	rs.task[uuid] = status
	return nil
}

func (rs *RamStorage) Post(uuid string, taskStatusResult string) error {
	if _, exists := rs.task[uuid]; exists {
		return errors.New("task already exists")
	}
	time.Sleep(15 * time.Second)
	taskStatusResult = strconv.Itoa(rand.Intn(100))
	rs.task[uuid] = taskStatusResult
	return nil
}

func (rs *RamStorage) Delete(uuid string) error {
	if _, exists := rs.task[uuid]; !exists {
		return errors.New("uuid not found")
	}
	delete(rs.task, uuid)
	return nil
}
