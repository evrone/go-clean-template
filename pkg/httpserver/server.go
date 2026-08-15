// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	App *fiber.App

	address         string
	prefork         bool
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	logger logger.Interface
}

// New -.
func New(l logger.Interface, opts ...Option) *Server {
	s := &Server{
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		logger:          l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	app := fiber.New(fiber.Config{
		Prefork:      s.prefork,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	s.App = app

	return s
}

// Start starts the HTTP server and blocks until context is canceled.
func (s *Server) Start(ctx context.Context) error {
	errChan := make(chan error, 1)

	// Start server in goroutine
	go func() {
		s.logger.Info("restapi server - Server - Started on %s", s.address)

		if err := s.App.Listen(s.address); err != nil {
			errChan <- err
		}
	}()

	// Wait for either error or context cancellation
	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)

	case <-ctx.Done():
		s.logger.Info("restapi server - Server - Shutting down...")

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.shutdownTimeout)
		defer cancel()

		if err := s.App.ShutdownWithContext(shutdownCtx); err != nil {
			s.logger.Error(err, "restapi server - Server - Shutdown error")

			return fmt.Errorf("shutdown error: %w", err)
		}

		s.logger.Info("restapi server - Server - Shutdown complete")

		return ctx.Err()
	}
}
