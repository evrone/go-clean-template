package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	*zap.Logger
}

func NewZapLogger(logLevel string) *ZapLogger {
	var (
		level  zapcore.Level
		logger *zap.Logger
	)

	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		log.Fatalf("zap init error: %s", err)
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("zap init error: %s", err)
	}

	return &ZapLogger{logger}
}

func (z *ZapLogger) Close() {
	_ = z.Logger.Sync() //nolint:errcheck // can't log it
}
