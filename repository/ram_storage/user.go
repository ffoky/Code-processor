package ram_storage

import (
	"github.com/google/uuid"
	"http_server/domain"
	"http_server/repository"
)

type User struct {
	users map[string]domain.User
}

func NewUser() *User {
	return &User{
		users: make(map[string]domain.User),
	}
}

func (rs *User) Get(uuid uuid.UUID) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (rs *User) Put(uuid uuid.UUID, Login string, Password string) error {
	//TODO implement me
	panic("implement me")
}

func (rs *User) Delete(uuid uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

//func (rs *User) Get(id uuid.UUID) (domain.User, error) {
//	user, exists := rs.users[id]
//	if !exists {
//		return domain.User{}, repository.NotFound
//	}
//	return user, nil
//}
//
//func (rs *User) Put(id uuid.UUID, login string, password string) error {
//	user, exists := rs.users[id]
//	if !exists {
//		return repository.NotFound
//	}
//
//	if login != "" {
//		user.Login = login
//	}
//
//	if password != "" {
//		user.Password = password
//	}
//
//	rs.users[id] = user
//	return nil
//}

func (rs *User) Post(id uuid.UUID, login string, password string) error {
	if _, exists := rs.users[login]; exists {
		return repository.ErrUserExists
	}

	rs.users[login] = domain.User{
		Uid:      id,
		Login:    login,
		Password: password,
	}
	return nil
}

func (rs *User) Login(login string, password string) (*domain.User, error) {
	user, exists := rs.users[login]
	if !exists {
		return nil, repository.ErrUserNotFound
	}

	if user.Password != password {
		return nil, repository.ErrUserNotFound
	}

	return &user, nil
}

//func (rs *User) Delete(id uuid.UUID) error {
//	if _, exists := rs.users[id]; !exists {
//		return repository.NotFound
//	}
//	delete(rs.users, id)
//	return nil
//}
