package api

import (
	"bytes"
	"encoding/json"
	"github.com/BrunoDM2943/church-members-api/internal/handler/api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpResult struct {
	status int
	body   []byte
}

type mapAssert func(parsedBody interface{})

func (result httpResult) assertStatus(t *testing.T, expected int) {
	assert.Equal(t, expected, result.status)
}

func (result httpResult) assert(t *testing.T, expected int, dto interface{}, assertBody mapAssert) {
	assert.Equal(t, expected, result.status)
	err := json.Unmarshal(result.body, dto)
	if err != nil {
		t.Error(err)
	} else {
		assertBody(dto)
	}
}

func newApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ApiErrorMiddleWare,
	})
	app.Use(recover.New())
	return app
}

func runTest(app *fiber.App, req *http.Request) httpResult {
	resp, err := app.Test(req, -1)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return httpResult{
		status: resp.StatusCode,
		body:   body,
	}
}

func buildPost(target, body string) *http.Request {
	return buildRequest("POST", target, body)
}

func buildPut(target, body string) *http.Request {
	return buildRequest("PUT", target, body)
}


func buildGet(target string) *http.Request {
	return buildRequest("GET", target, "")
}

func buildRequest(method, target, body string) *http.Request {
	req := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
