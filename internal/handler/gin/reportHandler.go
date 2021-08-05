package gin

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	report2 "github.com/BrunoDM2943/church-members-api/internal/modules/report"
	"net/http"

	"github.com/gin-gonic/gin"
)

//ReportHandler is a REST controller
type ReportHandler struct {
	reportGenerator report2.Service
}

//NewReportHandler builds a new ReportHandler
func NewReportHandler(reportGenerator report2.Service) *ReportHandler {
	return &ReportHandler{reportGenerator}
}

func (handler *ReportHandler) generateClassificationReport(c *gin.Context) {
	classification, err := new(enum.Classification).From(c.Param("classification"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{ Message: "Invalid classification: " + err.Error()})
		return
	}
	output, err := handler.reportGenerator.ClassificationReport(classification)
	handler.buildResponse(c, output, "classification_report.pdf", "application/pdf", err)
}

func (handler *ReportHandler) generateMarriageReport(c *gin.Context) {
	output, err := handler.reportGenerator.MarriageReport()
	handler.buildResponse(c, output, "marriage.csv", "application/csv", err)
}

func (handler *ReportHandler) generateBirthDayReport(c *gin.Context) {
	output, err := handler.reportGenerator.BirthdayReport()
	handler.buildResponse(c, output, "birthday.csv", "application/csv", err)
}

func (handler *ReportHandler) generateMembersReport(c *gin.Context) {
	output, err := handler.reportGenerator.MemberReport()
	handler.buildResponse(c, output, "members_report.pdf", "application/pdf", err)
}

func (handler *ReportHandler) generateLegalReport(c *gin.Context) {
	output, err := handler.reportGenerator.LegalReport()
	handler.buildResponse(c, output, "legal_report.pdf", "application/pdf", err)
}

func (handler *ReportHandler) buildResponse(c *gin.Context, data []byte, fileName string, contentType string, err error){
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{ Message: "Error generating report", Error: err})
	} else {
		c.Status(http.StatusOK)
		c.Header("Content-Type", contentType)
		c.Header("Content-Disposition", "attachment")
		c.Header("filename", fileName)
		c.Data(http.StatusOK, contentType, data)
	}
}