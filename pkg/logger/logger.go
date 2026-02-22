package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Interface -.
type Interface interface {
	Debug(message any, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message any, args ...any)
	Fatal(message any, args ...any)
}

// Logger -.
type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	logger := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &Logger{
		logger: new(logger),
	}
}

// Debug -.
func (l *Logger) Debug(message any, args ...any) {
	l.msg(zerolog.DebugLevel, message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...any) {
	l.log(zerolog.InfoLevel, message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...any) {
	l.log(zerolog.WarnLevel, message, args...)
}

// Error -.
func (l *Logger) Error(message any, args ...any) {
	if l.logger.GetLevel() == zerolog.DebugLevel {
		l.Debug(message, args...)
	}

	l.msg(zerolog.ErrorLevel, message, args...)
}

// Fatal -.
func (l *Logger) Fatal(message any, args ...any) {
	l.msg(zerolog.FatalLevel, message, args...)

	os.Exit(1)
}

func (l *Logger) log(level zerolog.Level, message string, args ...any) {
	if len(args) == 0 {
		l.logger.WithLevel(level).Msg(message)
	} else {
		l.logger.WithLevel(level).Msgf(message, args...)
	}
}

func (l *Logger) msg(level zerolog.Level, message any, args ...any) {
	switch msg := message.(type) {
	case error:
		l.log(level, msg.Error(), args...)
	case string:
		l.log(level, msg, args...)
	default:
		l.log(level, fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
