package main

import (
	"log"
	"os"
	"strconv"
)

type Conf struct {
	AppProbePort string
	PgURL        string
	PgPoolMax    int
	PgTableName  string
	RmqURL       string
	RmqQueueName string
}

func NewConfig() Conf {
	return Conf{
		AppProbePort: strEnv("APP_PROBE_PORT"),

		PgURL:       strEnv("PG_URL"),
		PgPoolMax:   intEnv("PG_POOL_MAX"),
		PgTableName: strEnv("PG_TABLE_NAME"),

		RmqURL:       strEnv("RMQ_URL"),
		RmqQueueName: strEnv("RMQ_QUEUE_NAME"),
	}
}

func strEnv(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok || len(value) == 0 {
		log.Fatalf("environment variable not declared: %s", env)
	}

	log.Println(env, "=", value)

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

	log.Println(env, "=", intValue)

	return intValue
}
