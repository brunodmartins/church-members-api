package api

import (
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/gofiber/fiber/v2"
)

func badRequest(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
		Message: "Invalid body received",
		Error:   err.Error(),
	})
}
