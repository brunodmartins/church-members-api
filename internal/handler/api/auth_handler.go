package api

import (
	"encoding/base64"
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	"net/http"
	"strings"

	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/BrunoDM2943/church-members-api/platform/security"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService security.Service
}

func NewAuthHandler(authService security.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler *AuthHandler) getToken(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" || len(handler.splitAuthHeader(authHeader)) != 2 {
		return handler.builderUnauthorizedError()
	}
	decodedHeader, err := handler.decodeHeader(handler.splitAuthHeader(authHeader)[1])
	if err != nil {
		return err
	}
	if !strings.Contains(decodedHeader, ":") {
		return handler.builderUnauthorizedError()
	}
	userName, password := handler.getUser(decodedHeader)
	token, err := handler.authService.GenerateToken(userName, password)
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusCreated).JSON(&dto.GetTokenResponse{Token: token})
}

func (handler *AuthHandler) getUser(header string) (string, string) {
	result := strings.Split(header, ":")
	return result[0], result[1]
}

func (handler *AuthHandler) decodeHeader(encodedHeader string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(encodedHeader)
	if err != nil {
		return "", handler.builderUnauthorizedError()
	}
	return string(result), nil
}

func (handler *AuthHandler) splitAuthHeader(authHeader string) []string {
	return strings.Split(authHeader, " ")
}

func (handler *AuthHandler) builderUnauthorizedError() apierrors.Error {
	return apierrors.NewApiError("Invalid Authorization header", http.StatusUnauthorized)
}
