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

func (handler *UserHandler) GetUserByName(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	user, err := handler.service.FindUser(ctx.UserContext(), name)
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(dto.NewGetUserResponse(user))
}

func (handler *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	requestBody := new(dto.UpdateUserRequest)
	_ = json.Unmarshal(ctx.Body(), &requestBody)
	if err := ValidateStruct(requestBody); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Invalid body received",
			Error:   err.Error(),
		})
	}
	user := requestBody.ToUser()
	user.UserName = name
	if err := handler.service.UpdateUser(ctx.UserContext(), user); err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusOK)
}
