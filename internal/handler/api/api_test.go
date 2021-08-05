package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

func newApp() *gin.Engine {
	return gin.Default()
}

func runTest(app *gin.Engine, req *http.Request) httpResult {
	writer := httptest.NewRecorder()
	app.ServeHTTP(writer, req)
	return httpResult{
		status: writer.Code,
		body:   writer.Body.Bytes(),
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
