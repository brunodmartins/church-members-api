package middleware

import (
	"net/http"
	"strings"

	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/brunodmartins/church-members-api/platform/security"
	"github.com/gofiber/fiber/v2"
)

const authHeader = "X-Auth-Token"

var routesWithAuthorization = map[string][]string{
	"/members":      {"viewMember", "editMember"},
	"/reports":      {"viewReport", "editReport"},
	"/churches":     {"viewChurch", "editChurch"},
	"/participants": {"viewParticipant", "editParticipant"},
	"/users":        {"viewUser", "editUser"},
	"/test":         {"viewTest", "editTest"},
}

var AuthMiddlewareMiddleWare = func(ctx *fiber.Ctx) error {

	token := ctx.Get(authHeader)
	if token == "" {
		return apierrors.NewApiError("Missing authorization token", http.StatusUnauthorized)
	}
	valid, claim := security.GetClaim(token)
	if !valid {
		return apierrors.NewApiError("Invalid authorization token", http.StatusForbidden)
	}

	route := detectRoute(ctx.Path())
	if route == "" {
		return apierrors.NewApiError("Route not found", http.StatusForbidden)
	}

	if ok := authorizeRoute(claim, route, ctx.Method()); !ok {
		return apierrors.NewApiError("User does not have required role", http.StatusForbidden)
	}

	ctx.SetUserContext(security.AddClaimToContext(claim, ctx.UserContext()))

	return ctx.Next()
}

func detectRoute(path string) string {
	for route := range routesWithAuthorization {
		if strings.HasPrefix(path, route) {
			return route
		}
	}
	return ""
}

func authorizeRoute(claim *security.Claim, route, method string) bool {
	allowedRoles := routesWithAuthorization[route]
	requiredPrefix := "view"
	if method != http.MethodGet {
		requiredPrefix = "edit"
	}

	var requiredRoles []string
	for _, r := range allowedRoles {
		if strings.HasPrefix(r, requiredPrefix) {
			requiredRoles = append(requiredRoles, r)
		}
	}

	if claim == nil || len(claim.Roles) == 0 {
		return false
	}

	for _, need := range requiredRoles {
		for _, userRole := range claim.Roles {
			if userRole == need {
				return true
			}
		}
	}

	return false
}
