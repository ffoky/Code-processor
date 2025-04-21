package domain

import (
	"sync"
)

type SessionProvider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

type SessionManager struct {
	CookieName  string
	Lock        sync.Mutex
	Provider    SessionProvider
	MaxLifetime int64
	Secure      bool
}

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}
