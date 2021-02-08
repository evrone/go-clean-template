package logger

import "github.com/rollbar/rollbar-go"

type RollbarLogger struct {
	*rollbar.Client
}

func NewRollbarLogger(token, env string) *RollbarLogger {
	client := rollbar.NewAsync(token, env, "", "", "")
	return &RollbarLogger{client}
}

func (r *RollbarLogger) Close() {
	r.Client.Wait()
}
