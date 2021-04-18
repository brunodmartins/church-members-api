package gin

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	"net/http"

	"github.com/BrunoDM2943/church-members-api/internal/service/report"

	"github.com/gin-gonic/gin"
)

//ReportHandler is a REST controller
type ReportHandler struct {
	reportGenerator report.Service
}

//NewReportHandler builds a new ReportHandler
func NewReportHandler(reportGenerator report.Service) *ReportHandler {
	return &ReportHandler{reportGenerator}
}

func isValidClassification(classification string) bool {
	return classification == "adult" || classification == "teen" || classification == "young" || classification == "children"
}

func (handler *ReportHandler) generateClassificationReport(c *gin.Context) {
	classification := c.Param("classification")
	if !isValidClassification(classification) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{ Message: "Invalid classification: " + classification})
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