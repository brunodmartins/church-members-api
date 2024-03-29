package api

import (
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
