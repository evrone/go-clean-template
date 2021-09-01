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

// ConnWaitTime -.
func ConnWaitTime(timeout time.Duration) Option {
	return func(s *Server) {
		s.conn.WaitTime = timeout
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(s *Server) {
		s.conn.Attempts = attempts
	}
}
