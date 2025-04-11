package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"http_server/api/http"
	_ "http_server/docs"
	pkgHttp "http_server/pkg/http"
	"http_server/repository/ram_storage"
	"http_server/usecases/service"
	_ "log"
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

	logrus.Infof("Starting server on %s", *addr)
	if err := pkgHttp.CreateAndRunServer(r, *addr); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
