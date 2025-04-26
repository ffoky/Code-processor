package http

import (
	"github.com/go-chi/chi/v5"
	googleId "github.com/google/uuid"
	"http_server/api/http/middleware"
	"http_server/api/http/types"
	_ "http_server/repository/ram_storage"
	"http_server/usecases"
	"http_server/usecases/service"
	"log"
	"net/http"
)

// Task represents an HTTP handler for managing task.
type Task struct {
	service usecases.Task
}

// NewTaskHandler creates a new instance of Task.
func NewTaskHandler(service usecases.Task) *Task {
	return &Task{service: service}
}

// @Summary Get task status
// @Description Get task status by its id
// @Tags task
// @Accept  json
// @Produce json
// @Param id path string true "Task tid" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200     {object}  types.GetTaskStatusHandlerResponse
// @Failure 400     {string}  string  "Bad request"
// @Failure 401     {string}  string  "Unauthorized"
// @Failure 404     {string}  string  "Task not found"
// @Security ApiKeyAuth
// @Router  /status/{id} [get]
func (t *Task) getTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetTaskStatusHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	task, err := t.service.Get(req)
	types.ProcessError(w, err, types.GetTaskStatusHandlerResponse{TaskStatus: task.Status}, http.StatusCreated)

}

// @Summary Get task result
// @Description Get task result by its id
// @Tags task
// @Accept  json
// @Produce json
// @Param id path string true "Task tid" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200 {object} types.GetTaskResultHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Task not found"
// @Security ApiKeyAuth
// @Router /result/{id} [get]
func (t *Task) getTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetTaskResultHandlerRequest(r)
	log.Printf("error: %v", err)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	task, err := t.service.Get(req)
	if err != nil {
		http.Error(w, "Error getting task result", http.StatusInternalServerError)
		return
	}
	types.ProcessError(w, err, types.GetTaskResultHandlerResponse{TaskResult: task.Result}, http.StatusOK)
}

// @Summary Create task
// @Description Create new task with the specified id and result
// @Tags task
// @Accept  json
// @Produce json
// @Param request body types.PostTaskHandlerRequest  true  "Task creation data"
// @Success 201 {object} types.PostTaskHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /task [post]
// @Security ApiKeyAuth
func (t *Task) postTaskHandler(w http.ResponseWriter, r *http.Request) {
	_, err := types.CreatePostTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	responseUUID, err := t.service.Post()
	types.ProcessError(w, err, &types.PostTaskHandlerResponse{TaskId: responseUUID}, http.StatusCreated)
}

// @Summary Delete task
// @Description Delete task by its id
// @Tags task
// @Accept  json
// @Produce json
// @Param id query types.DeleteTaskHandlerRequest true "Task tid" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200 {string} string "Task deleted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Task not found"
// @Router /task [delete]
// @Security ApiKeyAuth
func (t *Task) deleteHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateDeleteTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	taskId, _ := googleId.Parse(req.Id)
	err = t.service.Delete(taskId)
	types.ProcessError(w, err, nil, http.StatusOK)
}

// WithTaskHandlers registers task-related HTTP handlers.
func (t *Task) WithTaskHandlers(r chi.Router, sessionService *service.SessionService) {
	authMiddleware := middleware.AuthMiddleware(sessionService)

	r.Route("/task", func(r chi.Router) {
		r.With(authMiddleware).Post("/", t.postTaskHandler)
		r.With(authMiddleware).Delete("/", t.deleteHandler)
	})

	r.Route("/status", func(r chi.Router) {
		r.With(authMiddleware).Get("/{id}", t.getTaskStatusHandler)
	})

	r.Route("/result", func(r chi.Router) {
		r.With(authMiddleware).Get("/{id}", t.getTaskResultHandler)
	})
}
