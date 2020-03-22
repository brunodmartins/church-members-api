package handler

import (
	"net/http"

	"github.com/BrunoDM2943/church-members-api/reports"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportGenerator reports.ReportsGenerator
}

func NewReportHandler(reportGenerator reports.ReportsGenerator) ReportHandler {
	return ReportHandler{reportGenerator}
}

func (handler ReportHandler) SetUpRoutes(r *gin.Engine) {
	r.GET("/utils/members/aniversariantes", handler.GenerateBirthDayReport)
}

func (handler ReportHandler) GenerateBirthDayReport(c *gin.Context) {
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
