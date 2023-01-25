package middleware

import (
	"gitlab.com/vk-go/lectures-2022-2/pkg/session"
	"log"
)

type Middleware struct {
	key     []byte
	session session.SesManager
	logger  *log.Logger
}

func NewMiddleware(key []byte, session session.SesManager, logger *log.Logger) *Middleware {
	return &Middleware{key: key, session: session, logger: logger}
}
