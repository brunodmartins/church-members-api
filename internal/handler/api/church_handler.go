package api

import (
	"encoding/json"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/modules/church"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/gofiber/fiber/v2"
)

type ChurchHandler struct {
	service church.Service
}

func NewChurchHandler(service church.Service) *ChurchHandler {
	return &ChurchHandler{
		service: service,
	}
}

func (h *ChurchHandler) getStatistics(c *fiber.Ctx) error {
	stats, err := h.service.GetStatistics(c.UserContext(), c.Params("id"))
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(dto.ChurchStatisticsResponse{
		TotalMembers:                 stats.TotalMembers,
		AgeDistribution:              stats.AgeDistribution,
		TotalMembersByGender:         stats.TotalMembersByGender,
		TotalMembersByClassification: stats.TotalMembersByClassification,
	})
}

func (h *ChurchHandler) getChurchByID(c *fiber.Ctx) error {
	result, err := h.service.GetChurch(c.UserContext(), c.Params("id"))
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(dto.GetChurchResponse{
		ID:           result.ID,
		Name:         result.Name,
		Abbreviation: result.Abbreviation,
		Logo:         result.Logo,
	})
}

func (h *ChurchHandler) updateChurch(c *fiber.Ctx) error {
	id := c.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}

	request := new(dto.UpdateChurchRequest)
	if err := json.Unmarshal(c.Body(), request); err != nil {
		return badRequest(c, err)
	}
	if err := ValidateStruct(request); err != nil {
		return badRequest(c, err)
	}

	churchToUpdate := request.ToChurch()
	churchToUpdate.ID = id
	if err := h.service.UpdateChurch(c.UserContext(), churchToUpdate); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(dto.MessageResponse{Message: "Church updated successfully"})
}
