package grpcserver

import (
	"net"

	pbgrpc "google.golang.org/grpc"
)

// Option -.
type Option func(*Server)

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

// ServerOptions -.
func ServerOptions(opts ...pbgrpc.ServerOption) Option {
	return func(s *Server) {
		s.serverOpts = append(s.serverOpts, opts...)
	}
}
