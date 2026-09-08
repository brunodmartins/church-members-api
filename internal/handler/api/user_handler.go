package api

import (
	"encoding/json"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/gofiber/fiber/v2"
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

func (handler *UserHandler) SearchUsers(ctx *fiber.Ctx) error {
	users, err := handler.service.SearchUser(ctx.UserContext(), user.AllUsers())
	if err != nil {
		return err
	}
	result := make([]dto.GetUserResponse, len(users))
	for i, user := range users {
		result[i] = dto.NewGetUserResponse(user)
	}
	return ctx.Status(http.StatusOK).JSON(dto.SearchUsersResponse{
		Users: result,
	})
}
