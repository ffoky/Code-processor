package service

import (
	googleId "github.com/google/uuid"
	"http_server/api/http/types"
	"http_server/domain"
	"http_server/repository"
	"net/http"
	_ "net/http"
)

type User struct {
	repo    repository.User
	session *SessionService
}

func NewUser(repo repository.User, session *SessionService) *User {
	return &User{
		repo:    repo,
		session: session,
	}
}

func (s *User) Login(login, password string, w http.ResponseWriter, r *http.Request) (string, error) {
	user, err := s.repo.Login(login, password)
	if err != nil {
		return "", err
	}

	session, err := s.session.Start(w, r)
	if err != nil {
		return "", err
	}

	if err := session.Set("user_id", user.Uid); err != nil {
		return "", err
	}

	return session.SessionID(), nil
}

// TODO заменить парамер в Get
func (rs *User) Get(req types.GetTaskHandlerRequest) (domain.User, error) {
	return rs.repo.Get(req.Uuid)
}

func (rs *User) Put(uuid googleId.UUID, login string, password string) error {
	return rs.repo.Put(uuid, login, password)
}

func (rs *User) generateUUID() googleId.UUID {
	id := googleId.New()
	return id
}

func (rs *User) Post(login string, password string) error {
	uuid := rs.generateUUID()
	return rs.repo.Post(uuid, login, password)
}

func (rs *User) Delete(uuid googleId.UUID) error {
	return rs.repo.Delete(uuid)
}
