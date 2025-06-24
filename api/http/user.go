package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"http_server/api/http/types"
	"http_server/repository"
	"http_server/usecases"
	"net/http"
)

//TODO переделать регистрацию и логин, обработать ошибки

// User represents an HTTP handler for managing user.
type User struct {
	service usecases.User
}

// NewUserHandler creates a new instance of User.
func NewUserHandler(service usecases.User) *User {
	return &User{service: service}
}

// @Summary Create User
// @Description Create new user with the specified id, login and password
// @Tags user
// @Accept  json
// @Produce json
// @Param request body types.PostUserRegistrationHandlerRequest  true  "User creation data"
// @Success 201 {object} types.PostUserRegistrationHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Router /register [post]
func (u *User) postUserRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePostUserRegistrationHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = u.service.Post(req.Login, req.Password)
	if err != nil {
		logrus.WithError(err).Error("user already exists")
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	types.ProcessError(w, err, &types.PostUserRegistrationHandlerResponse{StatusCreated: http.StatusCreated}, http.StatusCreated)

}

// @Summary Login user
// @Description Login user by login and password
// @Tags user
// @Accept  json
// @Produce json
// @Param request body types.PostUserLoginHandlerRequest  true  "user login data"
// @Success 200 {object} types.PostUserLoginHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Router /login [post]
func (u *User) postUserLoginHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Start postUserLoginHandler")

	req, err := types.CreatePostUserLoginHandlerRequest(r)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse login request")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	sessionId, err := u.service.Login(req.Login, req.Password, w, r)
	if err != nil {
		if err == repository.ErrUserNotFound {
			logrus.WithError(err).Error("User doesnt exist")
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
	}
	types.ProcessError(w, err, &types.PostUserLoginHandlerResponse{Token: sessionId}, http.StatusOK)
}

// WithTaskHandlers registers user-related HTTP handlers.
func (u *User) WithUserHandlers(r chi.Router) {
	r.Route("/login", func(r chi.Router) {
		r.Post("/", u.postUserLoginHandler)
	})
	r.Route("/register", func(r chi.Router) {
		r.Post("/", u.postUserRegistrationHandler)
	})
}
