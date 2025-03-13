package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"http_server/storage"
	"net/http"
)

type Storage interface {
	Get(uuid uuid.UUID) (*string, error)
	Put(uuid uuid.UUID, result string, status string) error
	Post(uuid uuid.UUID, result string, status string) error
	Delete(uuid uuid.UUID) error
}

type Server struct {
	storage Storage
}

type getStatusResultResponseHandler struct {
	Uuid string `json:"status"`
}

func newServer(storage Storage) *Server {
	return &Server{storage: storage}
}

func (s *Server) getTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("key")
	if taskId == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	taskResult, err := s.storage.Get(taskId)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	taskStatus := "in_progress"

	_, _ = fmt.Fprintln(w, taskStatus)

}

func (s *Server) getTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("key")
	if uuid == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	value, err := s.storage.Get(uuid)
	if err != nil || value == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	_, _ = fmt.Fprintln(w, *value)
}

func (s *Server) putHandler(w http.ResponseWriter, r *http.Request) {
	var task map[string]string
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	key, okKey := task["key"]
	value, okValue := task["value"]

	if !okKey || !okValue {
		http.Error(w, "Missing key or value", http.StatusBadRequest)
		return
	}

	if err := s.storage.Put(key, value); err != nil {
		http.Error(w, "Failed to store value", http.StatusInternalServerError)
		return
	}
}

func generateUUID() uuid.UUID {
	id := uuid.New()
	return id
}

func UpdateTaskResult() {

}

func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	var task map[uuid.UUID]storage.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	responseUUID := generateUUID()
	result := "in_progress"
	_, _ = fmt.Fprintln(w, responseUUID)

	if err := s.storage.Post(responseUUID, status); err != nil {
		http.Error(w, "Failed to begin task", http.StatusInternalServerError)
		return
	}
}

func (s *Server) deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	if err := s.storage.Delete(key); err != nil {
		http.Error(w, "Failed to delete key", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CreateAndRunServer(storage Storage, addr string) error {
	server := newServer(storage)

	r := chi.NewRouter()

	r.Route("/task", func(r chi.Router) {
		r.Post("/", server.postHandler)
		r.Put("/", server.putHandler)
		r.Delete("/", server.deleteHandler)
	})
	r.Route("/status", func(r chi.Router) {
		r.Get("/", server.getTaskStatusHandler)
	})
	r.Route("/result", func(r chi.Router) {
		r.Get("/", server.getTaskResultHandler)
	})

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}
