package repository

import "http_server/domain"

type Provider interface {
	SessionInit(sid string) (domain.Session, error)
	SessionRead(sid string) (domain.Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}
