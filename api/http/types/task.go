package types

import (
	"encoding/json"
	"errors"
	"fmt"
	googleId "github.com/google/uuid"
	"http_server/repository"
	"net/http"
)

type GetTaskHandlerRequest struct {
	Uuid googleId.UUID `json:"id"`
}

func CreateGetTaskHandlerRequest(r *http.Request) (GetTaskHandlerRequest, error) {
	id := r.URL.Query().Get("id")
	if id == "" {
		return GetTaskHandlerRequest{}, fmt.Errorf("missing task id")
	}
	taskId, _ := googleId.Parse(id)
	return GetTaskHandlerRequest{Uuid: taskId}, nil
}

type GetTaskStatusHandlerResponse struct {
	TaskStatus string `json:"status"`
}

type GetTaskResultHandlerResponse struct {
	TaskResult string `json:"result"`
}

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

type PostTaskHandlerResponse struct {
	TaskId googleId.UUID `json:"taskId"`
}

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

func ProcessError(w http.ResponseWriter, err error, resp any) {
	if err != nil {
		if errors.Is(err, repository.NotFound) {
			http.Error(w, "Id not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		return
	}

	if resp != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintln(w, "Not found")
	}
}
