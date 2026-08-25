package api

import (
	"net/http"

	dto "github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/reportType"
	"github.com/brunodmartins/church-members-api/internal/modules/report"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/gofiber/fiber/v2"
)

// ReportHandler is a REST controller
type ReportHandler struct {
	reportGenerator report.Service
}

// NewReportHandler builds a new ReportHandler
func NewReportHandler(reportGenerator report.Service) *ReportHandler {
	return &ReportHandler{reportGenerator}
}

func (handler *ReportHandler) listReports(ctx *fiber.Ctx) error {
	reports := handler.reportGenerator.ListReports(ctx.UserContext())
	response := make([]dto.ReportResponse, 0, len(reports))
	for _, item := range reports {
		response = append(response, dto.ReportResponse{
			Name:         item.Name,
			Type:         string(item.Type),
			URL:          item.URL,
			CreationDate: item.CreationDate,
		})
	}
	return ctx.Status(http.StatusOK).JSON(response)
}

func (handler *ReportHandler) getURLForReport(ctx *fiber.Ctx) error {
	reportTypeName := reportType.Type(ctx.Params("reportType"))
	if !reportType.IsValidReport(reportTypeName) {
		return apierrors.NewApiError("Invalid report type", http.StatusBadRequest)
	}
	url, err := handler.reportGenerator.GetReport(ctx.UserContext(), reportTypeName)
	if err != nil {
		return err
	}
	ctx.Response().Header.Add("Location", url)
	return ctx.SendStatus(http.StatusTemporaryRedirect)
}

func (handler *ReportHandler) generateReport(ctx *fiber.Ctx) error {
	reportTypeName := reportType.Type(ctx.Params("reportType"))
	if !reportType.IsValidReport(reportTypeName) {
		return apierrors.NewApiError("Invalid report type", http.StatusBadRequest)
	}
	err := handler.reportGenerator.GenerateReport(ctx.UserContext(), reportTypeName)
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(dto.GenerateReportResponse{Message: "Report generated successfully", Type: string(reportTypeName)})
}
