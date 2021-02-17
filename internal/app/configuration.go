package app

import (
	"log"
	"os"
	"strconv"
)

type Conf struct {
	ServiceName        string
	ServiceVersion     string
	ZapLogLevel        string
	RollbarAccessToken string
	RollbarEnvironment string
	HTTPAPIPort        string
	PgURL              string
	PgPoolMax          int
	PgConnAttempts     int
	RmqURL             string
	RmqQueueName       string
}

func NewConfig() Conf {
	return Conf{
		ServiceName:        strEnv("SERVICE_NAME"),
		ServiceVersion:     strEnv("SERVICE_VERSION"),
		ZapLogLevel:        strEnv("ZAP_LOG_LEVEL"),
		RollbarAccessToken: strEnv("ROLLBAR_ACCESS_TOKEN"),
		RollbarEnvironment: strEnv("ROLLBAR_ENVIRONMENT"),
		HTTPAPIPort:        strEnv("HTTP_API_PORT"),
		PgURL:              strEnv("PG_URL"),
		PgPoolMax:          intEnv("PG_POOL_MAX"),
		PgConnAttempts:     intEnv("PG_CONN_ATTEMPTS"),
		RmqURL:             strEnv("RMQ_URL"),
		RmqQueueName:       strEnv("RMQ_QUEUE_NAME"),
	}
}

func strEnv(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok || len(value) == 0 {
		log.Fatalf("environment variable not declared: %s", env)
	}

	return value
}

func intEnv(env string) int {
	var intValue int

	value, ok := os.LookupEnv(env)
	if !ok || len(value) == 0 {
		log.Fatalf("environment variable not declared: %s", env)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("typecast error to integer: %s", err)
	}

	return intValue
}
