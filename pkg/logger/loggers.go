package logger

import (
	"fmt"

	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
)

type loggers struct {
	zap           *ZapLogger
	rollbar       *RollbarLogger
	defaultFields []Field
}

func (l *loggers) debug(msg string, fields ...Field) {
	fields = append(l.defaultFields, fields...)
	l.zap.Debug(msg, zapFields(fields)...)
}

func (l *loggers) info(msg string, fields ...Field) {
	fields = append(l.defaultFields, fields...)
	l.zap.Info(msg, zapFields(fields)...)
}

func (l *loggers) warn(msg string, fields ...Field) {
	fields = append(l.defaultFields, fields...)
	l.rollbar.MessageWithExtras(rollbar.WARN, msg, rollbarMap(fields))
	l.zap.Warn(msg, zapFields(fields)...)
}

func (l *loggers) error(err error, msg string, fields ...Field) {
	err = fmt.Errorf("%s: %w", msg, err)

	fields = append(l.defaultFields, fields...)
	l.rollbar.ErrorWithStackSkipWithExtras(rollbar.ERR, err, 3, rollbarMap(fields))
	l.zap.Error(err.Error(), zapFields(fields)...)
}

func (l *loggers) fatal(err error, msg string, fields ...Field) {
	err = fmt.Errorf("%s: %w", msg, err)

	fields = append(l.defaultFields, fields...)
	l.rollbar.ErrorWithStackSkipWithExtras(rollbar.CRIT, err, 3, rollbarMap(fields))
	l.rollbar.Close()
	l.zap.Fatal(err.Error(), zapFields(fields)...) // os.Exit()
}

func rollbarMap(fields []Field) map[string]interface{} {
	m := make(map[string]interface{}, len(fields)*2) //nolint:gomnd // fields number always 2
	for _, field := range fields {
		m[field.Key] = field.Val
	}

	return m
}

func zapFields(fields []Field) []zap.Field {
	s := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		s = append(s, zap.Reflect(field.Key, field.Val))
	}

	return s
}
