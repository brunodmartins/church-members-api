package middleware

import (
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/brunodmartins/church-members-api/platform/security"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

const authHeader = "X-Auth-Token"

var AuthMiddlewareMiddleWare = func(ctx *fiber.Ctx) error {

	token := ctx.Get(authHeader)
	if token == "" {
		return apierrors.NewApiError("Missing authorization token", http.StatusUnauthorized)
	}
	valid, claim := security.GetClaim(token)
	if !valid {
		return apierrors.NewApiError("Invalid authorization token", http.StatusForbidden)
	}
	ctx.SetUserContext(security.AddClaimToContext(claim, ctx.UserContext()))
	return ctx.Next()
}
