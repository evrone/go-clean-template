package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/evrone/go-service-template/internal/router"
)

type Server struct {
	server http.Server
	errors chan error
}

func NewServer(router *router.ProbeRouter, port string) *Server {
	return &Server{
		server: http.Server{
			Addr:    net.JoinHostPort("", port),
			Handler: router.ServeMux,
		},
		errors: make(chan error, 1),
	}
}

func (s *Server) Start() {
	go func() {
		s.errors <- s.server.ListenAndServe()
		close(s.errors)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.errors
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
