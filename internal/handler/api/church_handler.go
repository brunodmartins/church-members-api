package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/modules/church"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/gofiber/fiber/v2"
)

const maxChurchLogoSize = 5 * 1024 * 1024

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
		Language:     result.Language,
		Email:        result.Email,
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

func (h *ChurchHandler) updateChurchLogo(c *fiber.Ctx) error {
	id := c.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}

	fileHeader, err := c.FormFile("logo")
	if err != nil || fileHeader.Size >= maxChurchLogoSize {
		return apierrors.NewApiError("Logo must be a PNG file smaller than 5 MB", http.StatusBadRequest)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return apierrors.NewApiError("Unable to read logo file", http.StatusBadRequest)
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxChurchLogoSize))
	if err != nil {
		return apierrors.NewApiError("Unable to read logo file", http.StatusBadRequest)
	}
	if len(data) == 0 || len(data) >= maxChurchLogoSize || http.DetectContentType(data) != "image/png" {
		return apierrors.NewApiError("Logo must be a PNG file smaller than 5 MB", http.StatusBadRequest)
	}

	if err := h.service.UpdateChurchLogo(c.UserContext(), id, data); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(dto.MessageResponse{Message: "Church logo updated successfully"})
}
