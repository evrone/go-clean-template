package logger

var appLogger Logger //nolint:gochecknoglobals // it's necessary

func NewAppLogger(zap *ZapLogger, rollbar *RollbarLogger, serviceName, serviceVersion string) {
	fields := []Field{
		{"service-name", serviceName},
		{"service-version", serviceVersion},
	}
	appLogger = &loggers{zap, rollbar, fields}
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
