package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	mock_security "github.com/BrunoDM2943/church-members-api/platform/security/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAuthHandler_GetToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_security.NewMockService(ctrl)
	authHandler := NewAuthHandler(service)
	authHandler.SetUpRoutes(app)

	const userName = "user-test"
	const password = "password"

	t.Run("Success - 201", func(t *testing.T) {
		service.EXPECT().GenerateToken(userName, password).Return("token", nil)
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+encodeValue(buildHeaderValue(userName, password)))
		result := runTest(app, request)
		result.assertStatus(t, http.StatusCreated)
		assert.NotEmpty(t, result.cookies[0].Value)
	})
	t.Run("Fail - Error on service - 500", func(t *testing.T) {
		service.EXPECT().GenerateToken(userName, password).Return("", genericError)
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+encodeValue(buildHeaderValue(userName, password)))
		runTest(app, request).assertStatus(t, http.StatusInternalServerError)
	})
	t.Run("Fail - Header not encrypted", func(t *testing.T) {
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+ buildHeaderValue(userName, password))
		runTest(app, request).assertStatus(t, http.StatusUnauthorized)
	})
	t.Run("Fail - Header invalid", func(t *testing.T) {
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic xxx")
		runTest(app, request).assertStatus(t, http.StatusUnauthorized)
	})
	t.Run("Fail - Header invalid", func(t *testing.T) {
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic " + encodeValue("xxxx"))
		runTest(app, request).assertStatus(t, http.StatusUnauthorized)
	})
	t.Run("Fail - Header missing", func(t *testing.T) {
		request := buildGet("/users/token")
		runTest(app, request).assertStatus(t, http.StatusUnauthorized)
	})
}

func buildAuthorizationHeader(request *http.Request, value string) {
	request.Header.Set("Authorization", value)
}

func encodeValue(value string) string {
	return base64.StdEncoding.EncodeToString(bytes.NewBufferString(value).Bytes())
}

func buildHeaderValue(userName string, password string) string {
	return fmt.Sprintf("%s:%s", userName, password)
}
