package handler

import (
	"net/http"

	"github.com/BrunoDM2943/church-members-api/reports"
	"github.com/gin-gonic/gin"
)

type reportHandler struct {
	reportGenerator reports.ReportsGenerator
}

func NewReportHandler(reportGenerator reports.ReportsGenerator) reportHandler {
	return reportHandler{reportGenerator}
}

func (handler reportHandler) SetUpRoutes(r *gin.Engine) {
	r.GET("/reports/members", handler.generateMembersReport)
	r.GET("/reports/members/birthday", handler.generateBirthDayReport)
	r.GET("/reports/members/marriage", handler.generateMarriageReport)
	r.GET("/reports/members/legal", handler.generateLegalReport)

}

func (handler reportHandler) generateMarriageReport(c *gin.Context) {
	output, err := handler.reportGenerator.MariageReport()
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

func (handler reportHandler) generateBirthDayReport(c *gin.Context) {
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

func (handler reportHandler) generateMembersReport(c *gin.Context) {
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

func (handler reportHandler) generateLegalReport(c *gin.Context) {
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
