package api

import (
	"encoding/json"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/modules/student"
	"github.com/gofiber/fiber/v2"
)

type StudentHandler struct {
	service student.Service
}

func NewStudentHandler(service student.Service) *StudentHandler {
	return &StudentHandler{service: service}
}

func (handler *StudentHandler) postStudent(ctx *fiber.Ctx) error {
	request := new(dto.CreateStudentRequest)
	if err := json.Unmarshal(ctx.Body(), request); err != nil {
		return badRequest(ctx, err)
	}
	if err := ValidateStruct(request); err != nil {
		return badRequest(ctx, err)
	}

	id, err := handler.service.CreateStudent(ctx.UserContext(), request.ToStudent())
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusCreated).JSON(dto.CreateMemberResponse{ID: id})
}
