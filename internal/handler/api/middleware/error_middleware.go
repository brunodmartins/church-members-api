package middleware

import (
	"fmt"
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

var ApiErrorMiddleWare = func(c *fiber.Ctx, err error) error {
	log.WithError(err).Errorf("There was a error")
	if apiError, ok := err.(apierrors.Error); ok {
		return c.Status(apiError.StatusCode()).JSON(dto.ErrorResponse{
			Message: apiError.Error(),
			Error:   apiError,
		})
	} else {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: fmt.Sprintf("Unexpected error: %s", err.Error()),
			Error:   err,
		})
	}

}