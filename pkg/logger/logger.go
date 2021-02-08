package logger

import (
	"fmt"

	"github.com/evrone/go-service-template/internal/business-logic/domain"

	"github.com/rollbar/rollbar-go"

	"go.uber.org/zap"
)

type appLogger struct {
	zap     *ZapLogger
	rollbar *RollbarLogger
}

func NewAppLogger(zap *ZapLogger, rollbar *RollbarLogger) domain.Log {
	return &appLogger{zap, rollbar}
}

func (l *appLogger) Debug(msg string, fields ...domain.Field) {
	if fields == nil {
		l.rollbar.Message(rollbar.DEBUG, msg)
		l.zap.Debug(msg)
		return
	}

	l.rollbar.MessageWithExtras(rollbar.DEBUG, msg, rollbarMap(fields))
	l.zap.Debug(msg, zapFields(fields)...)
}

func (l *appLogger) Info(msg string, fields ...domain.Field) {
	if fields == nil {
		l.rollbar.Message(rollbar.INFO, msg)
		l.zap.Info(msg)
		return
	}

	l.rollbar.MessageWithExtras(rollbar.INFO, msg, rollbarMap(fields))
	l.zap.Info(msg, zapFields(fields)...)
}

func (l *appLogger) Warn(msg string, fields ...domain.Field) {
	if fields == nil {
		l.rollbar.Message(rollbar.WARN, msg)
		l.zap.Warn(msg)
		return
	}

	l.rollbar.MessageWithExtras(rollbar.WARN, msg, rollbarMap(fields))
	l.zap.Warn(msg, zapFields(fields)...)
}

func (l *appLogger) Error(err error, msg string, fields ...domain.Field) {
	err = fmt.Errorf("%s: %w", msg, err)

	if fields == nil {
		l.rollbar.ErrorWithStackSkip(rollbar.ERR, err, 4)
		l.zap.Error(err.Error())
		return
	}

	l.rollbar.ErrorWithStackSkipWithExtras(rollbar.ERR, err, 3, rollbarMap(fields))
	l.zap.Error(err.Error(), zapFields(fields)...)
}

func (l *appLogger) Fatal(err error, msg string, fields ...domain.Field) {
	err = fmt.Errorf("%s: %w", msg, err)

	if fields == nil {
		l.rollbar.ErrorWithStackSkip(rollbar.CRIT, err, 4)
		l.rollbar.Close()
		l.zap.Fatal(err.Error()) // os.Exit()
		return
	}

	l.rollbar.ErrorWithStackSkipWithExtras(rollbar.CRIT, err, 3, rollbarMap(fields))
	l.rollbar.Close()
	l.zap.Fatal(err.Error(), zapFields(fields)...) // os.Exit()
}

func rollbarMap(fields []domain.Field) map[string]interface{} {
	m := make(map[string]interface{}, len(fields)*2)
	for _, field := range fields {
		m[field.Key] = field.Val
	}
	return m
}

func zapFields(fields []domain.Field) []zap.Field {
	s := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		s = append(s, zap.Reflect(field.Key, field.Val))
	}
	return s
}
