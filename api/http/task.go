package http

import (
	"github.com/go-chi/chi/v5"
	googleId "github.com/google/uuid"
	"http_server/api/http/types"
	_ "http_server/repository/ram_storage"
	"http_server/usecases"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Task represents an HTTP handler for managing task.
type Task struct {
	service usecases.Task
}

// NewHandler creates a new instance of Task.
func NewHandler(service usecases.Task) *Task {
	return &Task{service: service}
}

// @Summary Get task status
// @Description Get task status by its id
// @Tags object
// @Accept  json
// @Produce json
// @Param id query string true "ID of the object"
// @Success 200 {object} types.GetTaskResultHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task not found"
// @Router /status [get]
func (t *Task) getTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	task, err := t.service.Get(*req)
	types.ProcessError(w, err, &types.GetTaskStatusHandlerResponse{TaskStatus: &task.Status})

}

// @Summary Get task status
// @Description Get task status by its id
// @Tags object
// @Accept  json
// @Produce json
// @Param id query string true "ID of the object"
// @Success 200 {object} types.GetTaskStatusHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task not found"
// @Router /result [get]
func (t *Task) getTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetTaskHandlerRequest(r)
	log.Printf("error: %v", err)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	task, err := t.service.Get(*req)
	types.ProcessError(w, err, &types.GetTaskResultHandlerResponse{TaskResult: &task.Result})

}

// @Summary Create or update task
// @Description Create or update task with the specified id,  status and result
// @Tags object
// @Accept  json
// @Produce json
// @Param id,status,result query string true "ID of the object"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad request"
// @Router /task [put]
func (t *Task) putHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePutTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	taskId, _ := googleId.Parse(req.Id)
	taskStatus := req.Status
	taskResult := req.Result
	err = t.service.Put(taskId, taskStatus, taskResult)
	types.ProcessError(w, err, http.StatusOK)
}

func (t *Task) generateUUID() googleId.UUID {
	id := googleId.New()
	return id
}

func (t *Task) completeTask(uuid googleId.UUID, w http.ResponseWriter) {
	time.Sleep(15 * time.Second)
	taskStatus := "ready"
	taskResult := strconv.Itoa(rand.Intn(100))
	if err := t.service.Put(uuid, taskStatus, taskResult); err != nil {
		http.Error(w, "Failed to complete task", http.StatusInternalServerError)
		return
	}
}

// @Summary Create task
// @Description Create new task with the specified id and result
// @Tags task
// @Accept  json
// @Produce json
// @Param key,value query string true "ID of the object"
// @Success 200 {string} string "TaskId"
// @Failure 400 {string} string "Bad request"
// @Router /task [post]
func (t *Task) postHandler(w http.ResponseWriter, r *http.Request) {
	_, err := types.CreatePostTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	responseUUID := t.generateUUID()
	go t.completeTask(responseUUID, w)
	taskStatus := "in_progress"

	err = t.service.Post(responseUUID, taskStatus, "")
	types.ProcessError(w, err, &types.PostTaskHandlerResponse{TaskId: responseUUID})
}

// @Summary Delete task
// @Description Delete task by its id
// @Tags object
// @Accept  json
// @Produce json
// @Param id query string true "ID of the object"
// @Success 200 {string} string "Task deleted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task not found"
// @Router /task [delete]
func (t *Task) deleteHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateDeleteTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	taskId, _ := googleId.Parse(req.Id)
	err = t.service.Delete(taskId)
	types.ProcessError(w, err, http.StatusOK)
}

// WithTaskHandlers registers task-related HTTP handlers.
func (t *Task) WithTaskHandlers(r chi.Router) {
	r.Route("/task", func(r chi.Router) {
		r.Post("/", t.postHandler)
		r.Put("/", t.putHandler)
		r.Delete("/", t.deleteHandler)
	})
	r.Route("/status", func(r chi.Router) {
		r.Get("/", t.getTaskStatusHandler)
	})
	r.Route("/result", func(r chi.Router) {
		r.Get("/", t.getTaskResultHandler)
	})
}
