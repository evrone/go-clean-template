package rest_api_test

import (
	"log"
	"net/http"
	"os"
	"testing"

	. "github.com/Eun/go-hit"
)

var basePath string

func TestMain(m *testing.M) {
	host, ok := os.LookupEnv("HOST")
	if !ok || len(host) == 0 {
		log.Fatalf("environment variable not declared: HOST")
	}

	basePath = "http://" + host + "/api/v1"

	code := m.Run()
	os.Exit(code)
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
			"original": "текст для перевода",
	}`
	Test(t,
		Description("DoTranslate Fail"),
		Post(basePath+"/translation/do-translate"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".error").Equal("invalid request body1"),
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
