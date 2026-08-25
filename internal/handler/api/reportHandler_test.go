package api

import (
	"net/http"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/reportType"
	"github.com/brunodmartins/church-members-api/internal/modules/report"
	mock_report "github.com/brunodmartins/church-members-api/internal/modules/report/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReportHandler_ListReports(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)
	expected := []report.Report{{
		Name:         "Member's report",
		Type:         reportType.MEMBER,
		URL:          "/reports/members",
		CreationDate: "21/08/2026 10:20:30",
	}}

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().ListReports(gomock.Any()).Return(expected)
		runTest(app, buildGet("/reports")).assert(t, http.StatusOK, new([]dto.ReportResponse), func(parsedBody interface{}) {
			response := parsedBody.(*[]dto.ReportResponse)
			assert.Equal(t, "Member's report", (*response)[0].Name)
			assert.Equal(t, "members", (*response)[0].Type)
			assert.Equal(t, "/reports/members", (*response)[0].URL)
			assert.Equal(t, "21/08/2026 10:20:30", (*response)[0].CreationDate)
		})
	})
}
func TestReportHandler_getURLForReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().GetReport(gomock.Any(), gomock.Eq(reportType.Type(reportType.MEMBER))).Return("url", nil)
		response := runTest(app, buildGet("/reports/members"))
		assert.Equal(t, http.StatusTemporaryRedirect, response.status)
		assert.Equal(t, "url", response.header.Get("Location"))
	})
	t.Run("Success - 400", func(t *testing.T) {
		runTest(app, buildGet("/reports/xxx")).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Success - 500", func(t *testing.T) {
		reports.EXPECT().GetReport(gomock.Any(), gomock.Eq(reportType.Type(reportType.MEMBER))).Return("", genericError)
		runTest(app, buildGet("/reports/members")).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestReportHandler_generateReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().GenerateReport(gomock.Any(), gomock.Eq(reportType.Type(reportType.MEMBER))).Return(nil)
		runTest(app, buildPost("/reports/members", nil)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 400", func(t *testing.T) {
		runTest(app, buildPost("/reports/invalid", nil)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().GenerateReport(gomock.Any(), gomock.Eq(reportType.Type(reportType.MEMBER))).Return(genericError)
		runTest(app, buildPost("/reports/members", nil)).assertStatus(t, http.StatusInternalServerError)
	})
}
