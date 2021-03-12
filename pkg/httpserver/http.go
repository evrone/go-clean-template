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
	notify chan error
}

func NewServer(handler http.Handler, port string) *Server {
	s := &Server{
		server: http.Server{
			Addr:    net.JoinHostPort("", port),
			Handler: handler,
		},
		notify: make(chan error, 1),
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint:gomnd // it's magic
	defer cancel()

	return s.server.Shutdown(ctx)
}
