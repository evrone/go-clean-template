package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/evrone/go-clean-template/pkg/logger"
	pbgrpc "google.golang.org/grpc"
)

const (
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	App *pbgrpc.Server

	address         string
	serverOpts      []pbgrpc.ServerOption
	shutdownTimeout time.Duration

	logger logger.Interface
}

// New -.
func New(l logger.Interface, opts ...Option) *Server {
	s := &Server{
		address:         _defaultAddr,
		shutdownTimeout: _defaultShutdownTimeout,
		logger:          l,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.App = pbgrpc.NewServer(s.serverOpts...)

	return s
}

// Start starts the gRPC server and blocks until context is canceled.
func (s *Server) Start(ctx context.Context) error {
	lc := net.ListenConfig{}

	ln, err := lc.Listen(ctx, "tcp", s.address)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	s.logger.Info("grpc server - Server - Started on %s", s.address)

	// Start graceful shutdown goroutine
	go func() {
		<-ctx.Done()
		s.logger.Info("grpc server - Server - Shutting down...")

		stopCh := make(chan struct{})

		go func() {
			s.App.GracefulStop()
			close(stopCh)
		}()

		select {
		case <-stopCh:
			s.logger.Info("grpc server - Server - Shutdown complete")
		case <-time.After(s.shutdownTimeout):
			s.logger.Info("grpc server - Server - Shutdown timeout, forcing stop")
			s.App.Stop()
		}
	}()

	// Serve blocks until server stops
	if err := s.App.Serve(ln); err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return ctx.Err()
}
