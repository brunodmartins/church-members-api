package api

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum/classification"
	report2 "github.com/BrunoDM2943/church-members-api/internal/modules/report"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

//ReportHandler is a REST controller
type ReportHandler struct {
	reportGenerator report2.Service
}

//NewReportHandler builds a new ReportHandler
func NewReportHandler(reportGenerator report2.Service) *ReportHandler {
	return &ReportHandler{reportGenerator}
}

func (handler *ReportHandler) generateClassificationReport(ctx *fiber.Ctx) error {
	classification, err := classification.From(ctx.Params("classification"))
	if err != nil {
		return apierrors.NewApiError("Invalid classification: " + err.Error(), http.StatusBadRequest)
	}
	output, err := handler.reportGenerator.ClassificationReport(classification)
	return handler.buildResponse(ctx, output, "classification_report.pdf", "application/pdf", err)
}

func (handler *ReportHandler) generateMarriageReport(ctx *fiber.Ctx) error {
	output, err := handler.reportGenerator.MarriageReport()
	return handler.buildResponse(ctx, output, "marriage.csv", "application/csv", err)
}

func (handler *ReportHandler) generateBirthDayReport(ctx *fiber.Ctx) error {
	output, err := handler.reportGenerator.BirthdayReport()
	return handler.buildResponse(ctx, output, "birthday.csv", "application/csv", err)
}

func (handler *ReportHandler) generateMembersReport(ctx *fiber.Ctx) error {
	output, err := handler.reportGenerator.MemberReport()
	return handler.buildResponse(ctx, output, "members_report.pdf", "application/pdf", err)
}

func (handler *ReportHandler) generateLegalReport(ctx *fiber.Ctx) error {
	output, err := handler.reportGenerator.LegalReport()
	return handler.buildResponse(ctx, output, "legal_report.pdf", "application/pdf", err)
}

func (handler *ReportHandler) buildResponse(ctx *fiber.Ctx, data []byte, fileName string, contentType string, err error) error{
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse{ Message: "Error generating report", Error: err})
	} else {
		ctx.Response().Header.Add("Content-Type", contentType)
		ctx.Response().Header.Add("Content-Disposition", "attachment")
		ctx.Response().Header.Add("filename", fileName)
		_, err = ctx.Status(http.StatusOK).Type(contentType).Write(data)
		return err
	}
}