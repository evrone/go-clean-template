package logger

type Logger interface {
	debug(msg string, fields ...Field)
	info(msg string, fields ...Field)
	warn(msg string, fields ...Field)
	error(err error, msg string, fields ...Field)
	fatal(err error, msg string, fields ...Field)
}

type Field struct {
	Key string
	Val interface{}
}
