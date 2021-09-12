package middleware

import (
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/BrunoDM2943/church-members-api/platform/security"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

const authHeader = "X-Auth-Token"

var AuthMiddlewareMiddleWare = func(ctx *fiber.Ctx) error {

	token := ctx.Get(authHeader)
	if token == "" {
		return apierrors.NewApiError("Missing authorization token", http.StatusUnauthorized)
	}
	if !security.IsValidToken(token) {
		return apierrors.NewApiError("Invalid authorization token", http.StatusForbidden)
	}
	return ctx.Next()
}
