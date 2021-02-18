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
}

func NewConfig() Conf {
	return Conf{
		ServiceName:        strEnv("GOT_SERVICE_NAME"),
		ServiceVersion:     strEnv("GOT_SERVICE_VERSION"),
		ZapLogLevel:        strEnv("GOT_ZAP_LOG_LEVEL"),
		RollbarAccessToken: strEnv("GOT_ROLLBAR_ACCESS_TOKEN"),
		RollbarEnvironment: strEnv("GOT_ROLLBAR_ENVIRONMENT"),
		HTTPAPIPort:        strEnv("GOT_HTTP_API_PORT"),
		PgURL:              strEnv("GOT_PG_URL"),
		PgPoolMax:          intEnv("GOT_PG_POOL_MAX"),
		PgConnAttempts:     intEnv("GOT_PG_CONN_ATTEMPTS"),
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
