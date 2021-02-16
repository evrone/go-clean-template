package rest_api_test

import (
	"net/http"
	"os"
	"testing"

	. "github.com/Eun/go-hit"
)

func TestMain(m *testing.M) {
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
		Post("http://localhost:8080/api/v1/translation/do-translate"),
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
		Post("http://localhost:8080/api/v1/translation/do-translate"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".error").Equal("invalid request body"),
	)
}

func TestHistory(t *testing.T) {
	Test(t,
		Description("History Success"),
		Get("http://localhost:8080/api/v1/translation/history"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"history":[{`),
	)
}
