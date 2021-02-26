package rest_api_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
)

var basePath string //nolint:gochecknoglobals // it's necessary

func TestMain(m *testing.M) {
	host, ok := os.LookupEnv("HOST")
	if !ok || host == "" {
		log.Fatalf("environment variable not declared: HOST")
	}

	err := healthCheck(host, 20)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	basePath = "http://" + host + "/api/v1"

	code := m.Run()
	os.Exit(code)
}

func healthCheck(host string, attempts int) error {
	var err error

	healthPath := "http://" + host + "/health"

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

func TestDoTranslate(t *testing.T) {
	body := `{
		"destination": "en",
		"original": "текст для перевода",
		"source": "auto"
	}`
	Test(t,
		Description("DoTranslate Success"),
		Post(basePath+"/translation/do-translate"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".translation").Equal("text for translation"),
	)

	body = `{
		"destination": "en",
		"original": "текст для перевода"
	}`
	Test(t,
		Description("DoTranslate Fail"),
		Post(basePath+"/translation/do-translate"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".error").Equal("invalid request body"),
	)
}

func TestHistory(t *testing.T) {
	Test(t,
		Description("History Success"),
		Get(basePath+"/translation/history"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"history":[{`),
	)
}
