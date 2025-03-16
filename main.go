package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "http_server/docs"
	"http_server/repository/ram_storage"
	"http_server/usecases/service"
	"log"

	"http_server/api/http"
	pkgHttp "http_server/pkg/http"
)

// @title Homework1
// @version 1.0
// @description This is homework server.

// @host localhost:8080
// @BasePath /
func main() {
	addr := flag.String("addr", ":8080", "address for http server")

	taskRepo := ram_storage.NewTask()
	taskService := service.NewTask(taskRepo)
	taskHandlers := http.NewHandler(taskService)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	taskHandlers.WithTaskHandlers(r)

	log.Printf("Starting server on %s", *addr)
	if err := pkgHttp.CreateAndRunServer(r, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
