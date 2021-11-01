package middleware

import (
	"fmt"
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

var ApiErrorMiddleWare = func(ctx *fiber.Ctx, err error) error {
	logrus.WithError(err).Errorf("Unexpeceted error")
	if apiError, ok := err.(apierrors.Error); ok {
		return ctx.Status(apiError.StatusCode()).JSON(dto.ErrorResponse{
			Message: apiError.Error(),
		})
	} else {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: fmt.Sprintf("Unexpected error: %s", err.Error()),
		})
	}
}
