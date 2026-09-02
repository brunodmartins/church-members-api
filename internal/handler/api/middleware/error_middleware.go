package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"

	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

var ApiErrorMiddleWare = func(ctx *fiber.Ctx, err error) error {
	if apiError, ok := err.(apierrors.Error); ok {
		return ctx.Status(apiError.StatusCode()).JSON(dto.ErrorResponse{
			Message: apiError.Error(),
		})
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		message := fiberError.Message
		if fiberError.Code == http.StatusNotFound {
			message = "Route not found"
		}
		return ctx.Status(fiberError.Code).JSON(dto.ErrorResponse{
			Message: message,
		})
	}

	logrus.WithError(err).Error("Unexpected Internal Server error")
	return ctx.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse{
		Message: fmt.Sprint("Unexpected Internal Server error"),
	})
}
