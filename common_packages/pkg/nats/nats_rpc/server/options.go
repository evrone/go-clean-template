package server

import "time"

// Option -.
type Option func(*Server)

// Timeout -.
func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}
