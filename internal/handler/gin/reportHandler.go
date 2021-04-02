package gin

import (
	"github.com/BrunoDM2943/church-members-api/internal/service/report"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportGenerator report.Service
}

func NewReportHandler(reportGenerator report.Service) *ReportHandler {
	return &ReportHandler{reportGenerator}
}

func (handler *ReportHandler) SetUpRoutes(r *gin.Engine) {
	r.GET("/reports/members/birthday", handler.generateBirthDayReport)
	r.GET("/reports/members/marriage", handler.generateMarriageReport)
	r.GET("/reports/members/legal", handler.generateLegalReport)
	r.GET("/reports/members/classification/:classification", handler.generateClassificationReport)
	r.GET("/reports/members", handler.generateMembersReport)

}

func isValidClassification(classification string) bool {
	if classification == "" {
		return false
	}
	if classification != "adult" && classification != "teen" && classification != "young" && classification != "children" {
		return false
	}
	return true
}

func (handler *ReportHandler) generateClassificationReport(c *gin.Context) {
	classification := c.Param("classification")
	if !isValidClassification(classification) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Invalid classification: " + classification})
	} else {
		output, err := handler.reportGenerator.ClassificationReport(classification)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Error generating report", "err": err.Error()})
		} else {
			c.Status(http.StatusOK)
			c.Header("Content-Type", "application/csv")
			c.Header("Content-Disposition", "attachment")
			c.Header("filename", "casamento.csv")
			c.Data(200, "application/csv", output)
		}
	}

}

func (handler *ReportHandler) generateMarriageReport(c *gin.Context) {
	output, err := handler.reportGenerator.MarriageReport()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Error generating report", "err": err.Error()})
	} else {
		c.Status(http.StatusOK)
		c.Header("Content-Type", "application/csv")
		c.Header("Content-Disposition", "attachment")
		c.Header("filename", "casamento.csv")
		c.Data(200, "application/csv", output)
	}
}

func (handler *ReportHandler) generateBirthDayReport(c *gin.Context) {
	output, err := handler.reportGenerator.BirthdayReport()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Error generating report", "err": err.Error()})
	} else {
		c.Status(http.StatusOK)
		c.Header("Content-Type", "application/csv")
		c.Header("Content-Disposition", "attachment")
		c.Header("filename", "aniversariantes.csv")
		c.Data(200, "application/csv", output)
	}
}

func (handler *ReportHandler) generateMembersReport(c *gin.Context) {
	output, err := handler.reportGenerator.MemberReport()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Error generating report", "err": err.Error()})
	} else {
		c.Status(http.StatusOK)
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment")
		c.Header("filename", "members.pdf")
		c.Data(200, "application/pdf", output)
	}
}

func (handler *ReportHandler) generateLegalReport(c *gin.Context) {
	output, err := handler.reportGenerator.LegalReport()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Error generating report", "err": err.Error()})
	} else {
		c.Status(http.StatusOK)
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment")
		c.Header("filename", "members_juridico.pdf")
		c.Data(200, "application/pdf", output)
	}
}
