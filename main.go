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
// @description Homework swagger api, added sessions and auth.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization: Bearer

// @host localhost:8080
// @BasePath /
func main() {

	addr := flag.String("addr", ":8080", "address for http server")

	sessionProvider := ram_storage.NewSessionProvider()
	sessionService := service.NewSessionService(
		sessionProvider,
		"gosessionid",
		3600, // 1 hour
	)
	taskRepo := ram_storage.NewTask()
	taskService := service.NewTask(taskRepo)
	taskHandlers := http.NewTaskHandler(taskService)
	userRepo := ram_storage.NewUser()
	userService := service.NewUser(userRepo, sessionService)
	userHandlers := http.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	taskHandlers.WithTaskHandlers(r, sessionService)
	userHandlers.WithUserHandlers(r)

	logrus.Infof("Starting server on %s", *addr)
	if err := pkgHttp.CreateAndRunServer(r, *addr); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
