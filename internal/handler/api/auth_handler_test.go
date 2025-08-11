package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/platform/security"
	mock_security "github.com/brunodmartins/church-members-api/platform/security/mock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
	var church = buildChurch(churchID)

	t.Run("Success - 201", func(t *testing.T) {
		service.EXPECT().IdentifyChurch(gomock.Eq(""), gomock.Eq(churchID)).Return(church, nil)
		service.EXPECT().GenerateToken(gomock.Eq(church), userName, password).Return("token", nil)
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+encodeValue(buildHeaderValue(userName, password)), churchID)
		runTest(app, request).assert(t, http.StatusCreated, &dto.GetTokenResponse{}, func(parsedBody interface{}) {
			assert.NotEmpty(t, parsedBody.(*dto.GetTokenResponse).Token)
		})

	})
	t.Run("Success - 201 - New login with abbreviation", func(t *testing.T) {
		service.EXPECT().IdentifyChurch(gomock.Eq(church.Abbreviation), gomock.Eq("")).Return(church, nil)
		service.EXPECT().GenerateToken(gomock.Eq(church), userName, password).Return("token", nil)
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+encodeValue(buildHeaderValue(userName, password)), "")
		request.Header.Set("x-church-abbreviation", church.Abbreviation)
		runTest(app, request).assert(t, http.StatusCreated, &dto.GetTokenResponse{}, func(parsedBody interface{}) {
			assert.NotEmpty(t, parsedBody.(*dto.GetTokenResponse).Token)
		})

	})
	t.Run("Fail - Error on service - 500", func(t *testing.T) {
		service.EXPECT().IdentifyChurch(gomock.Eq(""), gomock.Eq(churchID)).Return(church, nil)
		service.EXPECT().GenerateToken(gomock.Eq(church), userName, password).Return("", genericError)
		request := buildGet("/users/token")
		buildAuthorizationHeader(request, "Basic "+encodeValue(buildHeaderValue(userName, password)), churchID)
		runTest(app, request).assertStatus(t, http.StatusInternalServerError)
	})
	t.Run("Fail - church_id empty", func(t *testing.T) {
		service.EXPECT().IdentifyChurch(gomock.Eq(""), gomock.Eq("")).Return(nil, errors.New("generic error"))
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

func TestAuthHandler_ConfirmUserEmail(t *testing.T) {
	viper.Set("security.token.expiration", 1000)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_security.NewMockService(ctrl)
	authHandler := NewAuthHandler(service)
	authHandler.SetUpRoutes(app)

	username := "user-test"
	churchID := domain.NewID()
	user := &domain.User{
		UserName: username,
		Church: &domain.Church{
			ID: churchID,
		},
	}
	token := security.GenerateJWTToken(user)

	ctx := context.WithValue(context.Background(), "user", user)
	t.Run("Given a valid URL, when perform a GET operation, then return 200 OK", func(t *testing.T) {
		service.EXPECT().ConfirmEmail(gomock.Eq(ctx), gomock.Eq(username)).Return(nil)
		runTest(app, buildGet(fmt.Sprintf("/users/confirm?accessToken=%s", token))).assertStatus(t, http.StatusOK)
	})
	t.Run("Given a valid URL, when perform a GET operation, then return 500 error", func(t *testing.T) {
		service.EXPECT().ConfirmEmail(gomock.Eq(ctx), gomock.Eq(username)).Return(errors.New("generic error"))
		runTest(app, buildGet(fmt.Sprintf("/users/confirm?accessToken=%s", token))).assertStatus(t, http.StatusInternalServerError)
	})
	t.Run("Given an invalid URL, when perform a GET operation, then return error", func(t *testing.T) {
		runTest(app, buildGet(fmt.Sprint("/users/confirm"))).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Given an invalid URL, when perform a GET operation, then return error", func(t *testing.T) {
		runTest(app, buildGet(fmt.Sprint("/users/confirm?accessToken=abc"))).assertStatus(t, http.StatusForbidden)
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

func buildChurch(id string) *domain.Church {
	return &domain.Church{
		ID:           id,
		Name:         "test church",
		Abbreviation: "tc",
		Language:     "pt-br",
		Email:        "test@test.com",
	}
}
