// Package grpcserver implements HTTP server.
package grpcserver

import (
	"context"
	"errors"
	"net"

	"github.com/evrone/go-clean-template/pkg/logger"
	"golang.org/x/sync/errgroup"
	pbgrpc "google.golang.org/grpc"
)

const (
	_defaultAddr = ":80"
)

// Server -.
type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	App     *pbgrpc.Server
	notify  chan error
	address string

	logger logger.Interface
}

// New -.
func New(l logger.Interface, opts ...Option) *Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1) // Run only one goroutine

	s := &Server{
		ctx:     ctx,
		eg:      group,
		App:     pbgrpc.NewServer(),
		notify:  make(chan error, 1),
		address: _defaultAddr,
		logger:  l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start -.
func (s *Server) Start() {
	s.eg.Go(func() error {
		var lc net.ListenConfig

		ln, err := lc.Listen(s.ctx, "tcp", s.address)
		if err != nil {
			s.notify <- err

			close(s.notify)

			return err
		}

		err = s.App.Serve(ln)
		if err != nil {
			s.notify <- err

			close(s.notify)

			return err
		}

		return nil
	})

	s.logger.Info("grpc server - Server - Started")
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	var shutdownErrors []error

	s.App.GracefulStop() // Attention! Close connection first

	// Wait for all goroutines to finish and get any error
	err := s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "grpc server - Server - Shutdown - s.eg.Wait")

		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("grpc server - Server - Shutdown")

	return errors.Join(shutdownErrors...)
}
