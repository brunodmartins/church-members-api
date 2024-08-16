package api

import (
	"context"
	"encoding/json"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (handler *UserHandler) PostUser(ctx *fiber.Ctx) error {
	requestBody := new(dto.CreateUserRequest)
	_ = json.Unmarshal(ctx.Body(), &requestBody)
	if err := ValidateStruct(requestBody); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Invalid body received",
			Error:   err.Error(),
		})
	}
	if err := handler.service.SaveUser(ctx.UserContext(), requestBody.ToUser()); err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusCreated)
}

func (handler *UserHandler) ConfirmUserEmail(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	church := ctx.Query("church")

	if token == "" || church == "" {
		return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Invalid confirm email url provided",
		})
	}

	ctx.SetUserContext(context.WithValue(ctx.UserContext(), "church_id", church))
	if err := handler.service.ConfirmEmail(ctx.UserContext(), ctx.Params("user"), token); err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusOK)
}
