package api

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum/classification"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum/reportType"
	"github.com/BrunoDM2943/church-members-api/internal/modules/report"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// ReportHandler is a REST controller
type ReportHandler struct {
	reportGenerator report.Service
}

// NewReportHandler builds a new ReportHandler
func NewReportHandler(reportGenerator report.Service) *ReportHandler {
	return &ReportHandler{reportGenerator}
}

func (handler *ReportHandler) generateClassificationReport(ctx *fiber.Ctx) error {
	classification, err := classification.From(ctx.Params("classification"))
	if err != nil {
		return apierrors.NewApiError("Invalid classification: "+err.Error(), http.StatusBadRequest)
	}
	return handler.reportGenerator.ClassificationReport(ctx.UserContext(), classification)
}

func (handler *ReportHandler) generateMarriageReport(ctx *fiber.Ctx) error {
	return handler.reportGenerator.MarriageReport(ctx.UserContext())
}

func (handler *ReportHandler) generateBirthDayReport(ctx *fiber.Ctx) error {
	return handler.reportGenerator.BirthdayReport(ctx.UserContext())
}

func (handler *ReportHandler) generateMembersReport(ctx *fiber.Ctx) error {
	return handler.reportGenerator.MemberReport(ctx.UserContext())
}

func (handler *ReportHandler) generateLegalReport(ctx *fiber.Ctx) error {
	return handler.reportGenerator.LegalReport(ctx.UserContext())
}

func (handler *ReportHandler) getURLForReport(ctx *fiber.Ctx) error {
	reportTypeName := ctx.Params("reportType")
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
