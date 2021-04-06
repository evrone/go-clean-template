package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var appLogger Logger //nolint:gochecknoglobals // it's ok

func NewAppLogger(zapLogger *ZapLogger, appName, appVersion string) {
	fields := []Field{
		{"app-name", appName},
		{"app-version", appVersion},
	}
	appLogger = &logger{zapLogger, fields}
}

func Debug(msg string, fields ...Field) {
	appLogger.debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	appLogger.info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	appLogger.warn(msg, fields...)
}

func Error(err error, msg string, fields ...Field) {
	appLogger.error(err, msg, fields...)
}

func Fatal(err error, msg string, fields ...Field) {
	appLogger.fatal(err, msg, fields...)
}

type logger struct {
	zap           *ZapLogger
	defaultFields []Field
}

func (l *logger) debug(msg string, fields ...Field) {
	fields = append(l.defaultFields, fields...)
	l.zap.Debug(msg, zapFields(fields)...)
}

func (l *logger) info(msg string, fields ...Field) {
	fields = append(l.defaultFields, fields...)
	l.zap.Info(msg, zapFields(fields)...)
}

func (l *logger) warn(msg string, fields ...Field) {
	fields = append(l.defaultFields, fields...)
	l.zap.Warn(msg, zapFields(fields)...)
}

func (l *logger) error(err error, msg string, fields ...Field) {
	err = fmt.Errorf("%s: %w", msg, err)

	fields = append(l.defaultFields, fields...)
	l.zap.Error(err.Error(), zapFields(fields)...)
}

func (l *logger) fatal(err error, msg string, fields ...Field) {
	err = fmt.Errorf("%s: %w", msg, err)

	fields = append(l.defaultFields, fields...)
	l.zap.Fatal(err.Error(), zapFields(fields)...) // os.Exit()
}

func zapFields(fields []Field) []zap.Field {
	s := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		s = append(s, zap.Reflect(field.Key, field.Val))
	}

	return s
}
