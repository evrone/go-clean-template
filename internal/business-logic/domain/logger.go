package domain

var Logger Log

type Log interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(err error, msg string, fields ...Field)
	Fatal(err error, msg string, fields ...Field)
}

type Field struct {
	Key string
	Val interface{}
}
