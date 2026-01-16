package api

import (
	"encoding/json"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/modules/participant"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/gofiber/fiber/v2"
)

type ParticipantHandler struct {
	service participant.Service
}

func NewParticipantHandler(service participant.Service) *ParticipantHandler {
	return &ParticipantHandler{service: service}
}

func (h *ParticipantHandler) postParticipant(ctx *fiber.Ctx) error {
	req := new(dto.CreateParticipantRequest)
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return h.badRequest(ctx, err)
	}
	if err := ValidateStruct(req); err != nil {
		return h.badRequest(ctx, err)
	}
	id, err := h.service.CreateParticipant(ctx.UserContext(), req.ToParticipant())
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusCreated).JSON(dto.CreateMemberResponse{ID: id})
}

func (h *ParticipantHandler) badRequest(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid body received", Error: err.Error()})
}

func (h *ParticipantHandler) getParticipant(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	p, err := h.service.GetParticipant(ctx.UserContext(), id)
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(dto.NewGetParticipantResponse(p))
}

func (h *ParticipantHandler) searchParticipant(ctx *fiber.Ctx) error {
	queryFilters := participant.QueryBuilder{}
	if name := ctx.Query("name"); name != "" {
		queryFilters.AddFilter("name", name)
	}
	if gender := ctx.Query("gender"); gender != "" {
		queryFilters.AddFilter("gender", gender)
	}
	participants, err := h.service.SearchParticipant(ctx.UserContext(), queryFilters.ToSpecification())
	if err != nil {
		return err
	}
	resp := make([]*dto.GetParticipantResponse, 0)
	for _, p := range participants {
		resp = append(resp, dto.NewGetParticipantResponse(p))
	}
	return ctx.Status(http.StatusOK).JSON(resp)
}

func (h *ParticipantHandler) updateParticipant(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	req := new(dto.CreateParticipantRequest)
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return h.badRequest(ctx, err)
	}
	p := req.ToParticipant()
	p.ID = id
	if err := h.service.UpdateParticipant(ctx.UserContext(), p); err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(dto.MessageResponse{Message: "Participant updated successfully"})
}

func (h *ParticipantHandler) deleteParticipant(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	if err := h.service.DeleteParticipant(ctx.UserContext(), id); err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(dto.MessageResponse{Message: "Participant deleted successfully"})
}
