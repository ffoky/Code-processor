package types

import (
	"encoding/json"
	"fmt"
	googleId "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
)

// GetTaskHandlerRequest represents request for getting task status
// swagger:param parameters GetTaskHandlerRequest
type GetTaskHandlerRequest struct {
	Uuid googleId.UUID `json:"id"`
}

func CreateGetTaskStatusHandlerRequest(r *http.Request) (GetTaskHandlerRequest, error) {
	id := path.Base(r.URL.Path)

	taskId, err := googleId.Parse(id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":     id,
			"error":  err.Error(),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Error("Failed to parse task tid")
		return GetTaskHandlerRequest{}, fmt.Errorf("invalid task id format")
	}

	return GetTaskHandlerRequest{Uuid: taskId}, nil
}

func CreateGetTaskResultHandlerRequest(r *http.Request) (GetTaskHandlerRequest, error) {
	id := path.Base(r.URL.Path)
	if id == "" {
		return GetTaskHandlerRequest{}, fmt.Errorf("missing task id")
	}
	taskId, _ := googleId.Parse(id)
	return GetTaskHandlerRequest{Uuid: taskId}, nil
}

// GetTaskStatusHandlerResponse represents response with task status
// swagger:response getTaskStatusResponse
type GetTaskStatusHandlerResponse struct {
	TaskStatus string `json:"status"`
}

// GetTaskResultHandlerResponse represents response with task result
// swagger:response types.getTaskResultResponse
type GetTaskResultHandlerResponse struct {
	TaskResult string `json:"result"`
}

// PostTaskHandlerRequest represents task post request
// swagger:param parameters PostHandler
type PostTaskHandlerRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func CreatePostTaskHandlerRequest(r *http.Request) (*PostTaskHandlerRequest, error) {
	var req PostTaskHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

// PostTaskHandlerResponse represents response with created task tid
// swagger:response postTaskResponse
type PostTaskHandlerResponse struct {
	TaskId googleId.UUID `json:"taskId"`
}

// DeleteTaskHandlerRequest represents request for deleting task
// swagger:param parameters DeleteTaskHandlerRequest
type DeleteTaskHandlerRequest struct {
	Id string `json:"id"`
}

func CreateDeleteTaskHandlerRequest(r *http.Request) (*DeleteTaskHandlerRequest, error) {
	id := r.URL.Query().Get("id")
	if id == "" {
		return nil, fmt.Errorf("missing id")
	}
	return &DeleteTaskHandlerRequest{Id: id}, nil
}
