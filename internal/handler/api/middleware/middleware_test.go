package middleware

import (
	"errors"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/brunodmartins/church-members-api/platform/security"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
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
	t.Run("Header  - 200", func(t *testing.T) {
		req := buildRequest()
		req.Header.Set(authHeader, buildToken())
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
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

func buildToken() string {
	return security.GenerateJWTToken(&domain.User{
		ID:       "id",
		UserName: "test",
	})
}

func buildRequest() *http.Request {
	req := httptest.NewRequest("GET", "/test", nil)
	return req
}
