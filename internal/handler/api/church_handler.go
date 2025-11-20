package api

import (
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/modules/church"
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
