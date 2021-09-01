package rmqrpc

import "errors"

var (
	ErrTimeout        = errors.New("timeout")
	ErrInternalServer = errors.New("internal server error")
	ErrBadHandler     = errors.New("unregistered handler")
)

const (
	Success = "success"
)
