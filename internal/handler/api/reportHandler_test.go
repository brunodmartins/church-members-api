package api

import (
	"github.com/brunodmartins/church-members-api/internal/constants/enum/classification"
	mock_report "github.com/brunodmartins/church-members-api/internal/modules/report/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

func TestBirthDayReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().BirthdayReport(gomock.Any()).Return(nil)
		runTest(app, buildPost("/reports/members/birthday", nil)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().BirthdayReport(gomock.Any()).Return(genericError)
		runTest(app, buildPost("/reports/members/birthday", nil)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestMarriageReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().MarriageReport(gomock.Any()).Return(nil)
		runTest(app, buildPost("/reports/members/marriage", nil)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().MarriageReport(gomock.Any()).Return(genericError)
		runTest(app, buildPost("/reports/members/marriage", nil)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestLegalReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().LegalReport(gomock.Any()).Return(nil)
		runTest(app, buildPost("/reports/members/legal", nil)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().LegalReport(gomock.Any()).Return(genericError)
		runTest(app, buildPost("/reports/members/legal", nil)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestMemberReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().MemberReport(gomock.Any()).Return(nil)
		runTest(app, buildPost("/reports/members", nil)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().MemberReport(gomock.Any()).Return(genericError)
		runTest(app, buildPost("/reports/members", nil)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestClassificationReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200 - CHILDREN", func(t *testing.T) {
		reports.EXPECT().ClassificationReport(gomock.Any(), gomock.Eq(classification.CHILDREN))
		runTest(app, buildPost("/reports/members/classification/children", nil)).assertStatus(t, http.StatusOK)
	})
	t.Run("Success - 200 - TEEN", func(t *testing.T) {
		reports.EXPECT().ClassificationReport(gomock.Any(), gomock.Eq(classification.TEEN))
		runTest(app, buildPost("/reports/members/classification/teen", nil)).assertStatus(t, http.StatusOK)
	})
	t.Run("Success - 200 - ADULT", func(t *testing.T) {
		reports.EXPECT().ClassificationReport(gomock.Any(), gomock.Eq(classification.ADULT))
		runTest(app, buildPost("/reports/members/classification/adult", nil)).assertStatus(t, http.StatusOK)
	})
	t.Run("Success - 200 - YOUNG", func(t *testing.T) {
		reports.EXPECT().ClassificationReport(gomock.Any(), gomock.Eq(classification.YOUNG))
		runTest(app, buildPost("/reports/members/classification/young", nil)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 400", func(t *testing.T) {
		runTest(app, buildPost("/reports/members/classification/X", nil)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().ClassificationReport(gomock.Any(), gomock.Eq(classification.YOUNG)).Return(genericError)
		runTest(app, buildPost("/reports/members/classification/young", nil)).assertStatus(t, http.StatusInternalServerError)
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
		reports.EXPECT().GetReport(gomock.Any(), gomock.Eq("members")).Return("url", nil)
		response := runTest(app, buildGet("/reports/members"))
		assert.Equal(t, http.StatusTemporaryRedirect, response.status)
		assert.Equal(t, "url", response.header.Get("Location"))
	})
	t.Run("Success - 400", func(t *testing.T) {
		runTest(app, buildGet("/reports/xxx")).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Success - 500", func(t *testing.T) {
		reports.EXPECT().GetReport(gomock.Any(), gomock.Eq("members")).Return("", genericError)
		runTest(app, buildGet("/reports/members")).assertStatus(t, http.StatusInternalServerError)
	})
}
