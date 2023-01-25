package server

import (
	"gitlab.com/vk-go/lectures-2022-2/pkg/service"
	"log"
)

type Server struct {
	service *service.Service
	log     *log.Logger
}

func NewServer(service *service.Service, log *log.Logger) *Server {
	return &Server{
		service: service,
		log:     log,
	}
}
