package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/handler/api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpResult struct {
	status  int
	body    []byte
	cookies []*http.Cookie
	header  http.Header
}

var emptyJson = []byte("{}")
var badJson = []byte("")

type mapAssert func(parsedBody interface{})

func (result httpResult) assertStatus(t *testing.T, expected int) {
	assert.Equal(t, expected, result.status, string(result.body))
}

func (result httpResult) assert(t *testing.T, expected int, dto interface{}, assertBody mapAssert) {
	assert.Equal(t, expected, result.status, string(result.body))
	err := json.Unmarshal(result.body, dto)
	if err != nil {
		t.Error(err)
	} else {
		assertBody(dto)
	}
}

// assertJSON asserts the HTTP status and that the response body is JSON-equivalent to the marshaled expectedBody,
// guaranteeing full compliance of the API response shape (fields, names and structure).
func (result httpResult) assertJSON(t *testing.T, expected int, expectedBody interface{}) {
	assert.Equal(t, expected, result.status, string(result.body))
	expectedJSON, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatal(err)
	}
	assert.JSONEq(t, string(expectedJSON), string(result.body))
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
		status:  resp.StatusCode,
		header:  resp.Header,
		body:    body,
		cookies: resp.Cookies(),
	}
}

func buildPost(target string, body []byte) *http.Request {
	return buildRequest("POST", target, body)
}

func buildDelete(target string, body []byte) *http.Request {
	return buildRequest("DELETE", target, body)
}

func buildPut(target string, body []byte) *http.Request {
	return buildRequest("PUT", target, body)
}

func buildGet(target string) *http.Request {
	return buildRequest("GET", target, []byte(""))
}

func buildRequest(method, target string, body []byte) *http.Request {
	req := httptest.NewRequest(method, target, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func getMock(filename string) []byte {
	result, _ := ioutil.ReadFile(fmt.Sprintf("./resources/%s", filename))
	return result
}
