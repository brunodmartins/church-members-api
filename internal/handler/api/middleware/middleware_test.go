package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/brunodmartins/church-members-api/platform/security"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	viper.Set("security.token.expiration", 1)
	app := fiber.New(fiber.Config{
		ErrorHandler: ApiErrorMiddleWare,
	})
	app.Use(AuthMiddlewareMiddleWare)
	app.Use(func(ctx *fiber.Ctx) error {
		assert.NotNil(t, ctx.UserContext().Value("user"))
		return ctx.Next()
	})
	app.Get("/test", func(ctx *fiber.Ctx) error {
		return nil
	})
	app.Post("/test/x", func(ctx *fiber.Ctx) error {
		return nil
	})
	t.Run("Success - User with viewTest role  - 200", func(t *testing.T) {
		req := buildRequest()
		req.Header.Set(authHeader, buildToken("viewTest"))
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("Fail - User with editTest role  - 403", func(t *testing.T) {
		req := buildRequest()
		req.Header.Set(authHeader, buildToken("editTest"))
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
	t.Run("Success - User with editTest role  - 200", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test/x", nil)
		req.Header.Set(authHeader, buildToken("editTest"))
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("Fail - User with editTest role  - 403", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test/x", nil)
		req.Header.Set(authHeader, buildToken("viewTest"))
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
	t.Run("Missing Header - 401", func(t *testing.T) {
		req := buildRequest()
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
	t.Run("Header  - 403", func(t *testing.T) {
		req := buildRequest()
		req.Header.Set(authHeader, "invalid-token")
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func TestErrorMiddleware(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ApiErrorMiddleWare,
	})
	var testError error
	app.Get("/test", func(ctx *fiber.Ctx) error {
		return testError
	})
	t.Run("Treat API error", func(t *testing.T) {
		testError = apierrors.NewApiError("This is an error", http.StatusNotFound)
		req := buildRequest()
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("Treat generic error", func(t *testing.T) {
		testError = errors.New("generic error")
		req := buildRequest()
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestI18NMiddleware(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ApiErrorMiddleWare,
	})
	app.Use(I18NMiddleware)
	app.Get("/test", func(ctx *fiber.Ctx) error {
		assert.NotNil(t, ctx.Get("i18n"))
		return nil
	})
	t.Run("With Header  - 200", func(t *testing.T) {
		req := buildRequest()
		req.Header.Set(languageHeader, "pt-BR")
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Without Header  - 200", func(t *testing.T) {
		req := buildRequest()
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func buildToken(roles ...string) string {
	return security.GenerateJWTToken(&domain.User{
		ID:       "id",
		UserName: "test",
		Roles:    roles,
	})
}

func buildRequest() *http.Request {
	req := httptest.NewRequest("GET", "/test", nil)
	return req
}
