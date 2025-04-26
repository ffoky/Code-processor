package repository

import "http_server/domain"

type ObjectSender interface {
	Send(object domain.Task)
}
