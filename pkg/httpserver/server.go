// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	notify chan error
}

func NewServer(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler: handler,
		// Default
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         ":80",
	}

	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
	}

	// Set options
	for _, opt := range opts {
		opt(s)
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //nolint:gomnd // it's magic
	defer cancel()

	return s.server.Shutdown(ctx)
}
