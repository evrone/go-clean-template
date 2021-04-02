package rmqrpc

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrTimeout        = Error("timeout")
	ErrInternalServer = Error("internal server error")
	ErrBadHandler     = Error("unregistered handler")
)

const (
	Success = "success"
)
