// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Server struct {
	server http.Server
	errors chan error
}

func NewServer(handler http.Handler, port string) *Server {
	return &Server{
		server: http.Server{
			Addr:    net.JoinHostPort("", port),
			Handler: handler,
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint:gomnd // it's magic
	defer cancel()

	return s.server.Shutdown(ctx)
}
