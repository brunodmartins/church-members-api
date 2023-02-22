package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	mock_security "github.com/brunodmartins/church-members-api/platform/security/mock"
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
	var churchID = domain.NewID()

	t.Run("Success - 201", func(t *testing.T) {
		service.EXPECT().GenerateToken(churchID, userName, password).Return("token", nil)
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+encodeValue(buildHeaderValue(userName, password)), churchID)
		runTest(app, request).assert(t, http.StatusCreated, &dto.GetTokenResponse{}, func(parsedBody interface{}) {
			assert.NotEmpty(t, parsedBody.(*dto.GetTokenResponse).Token)
		})

	})
	t.Run("Fail - Error on service - 500", func(t *testing.T) {
		service.EXPECT().GenerateToken(churchID, userName, password).Return("", genericError)
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+encodeValue(buildHeaderValue(userName, password)), churchID)
		runTest(app, request).assertStatus(t, http.StatusInternalServerError)
	})
	t.Run("Fail - church_id empty", func(t *testing.T) {
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+encodeValue(buildHeaderValue(userName, password)), domain.NewID())
		request.Header.Del("church_id")
		runTest(app, request).assertStatus(t, http.StatusUnauthorized)
	})
	t.Run("Fail - Header not encrypted", func(t *testing.T) {
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+buildHeaderValue(userName, password), domain.NewID())
		runTest(app, request).assertStatus(t, http.StatusUnauthorized)
	})
	t.Run("Fail - Header invalid", func(t *testing.T) {
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic xxx", domain.NewID())
		runTest(app, request).assertStatus(t, http.StatusUnauthorized)
	})
	t.Run("Fail - Header invalid", func(t *testing.T) {
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+encodeValue("xxxx"), domain.NewID())
		runTest(app, request).assertStatus(t, http.StatusUnauthorized)
	})
	t.Run("Fail - Header missing", func(t *testing.T) {
		request := buildGet("/users/token")
		runTest(app, request).assertStatus(t, http.StatusUnauthorized)
	})
}

func buildAuthorizationHeader(request *http.Request, auth string, churchID string) {
	request.Header.Set("Authorization", auth)
	request.Header.Set("church_id", churchID)
}

func encodeValue(value string) string {
	return base64.StdEncoding.EncodeToString(bytes.NewBufferString(value).Bytes())
}

func buildHeaderValue(userName string, password string) string {
	return fmt.Sprintf("%s:%s", userName, password)
}
