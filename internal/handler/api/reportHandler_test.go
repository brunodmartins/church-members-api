package api

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	mock_report2 "github.com/BrunoDM2943/church-members-api/internal/modules/report/mock"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestBirthDayReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report2.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().BirthdayReport().Return([]byte{}, nil)
		runTest(app, buildGet("/reports/members/birthday")).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().BirthdayReport().Return([]byte{}, genericError)
		runTest(app, buildGet("/reports/members/birthday")).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestMarriageReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report2.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().MarriageReport().Return([]byte{}, nil)
		runTest(app, buildGet("/reports/members/marriage")).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().MarriageReport().Return([]byte{}, genericError)
		runTest(app, buildGet("/reports/members/marriage")).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestLegalReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report2.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().LegalReport().Return([]byte{}, nil)
		runTest(app, buildGet("/reports/members/legal")).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().LegalReport().Return([]byte{}, genericError)
		runTest(app, buildGet("/reports/members/legal")).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestMemberReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report2.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		reports.EXPECT().MemberReport().Return([]byte{}, nil)
		runTest(app, buildGet("/reports/members")).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().MemberReport().Return([]byte{}, genericError)
		runTest(app, buildGet("/reports/members")).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestClassificationReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()
	reports := mock_report2.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(app)

	t.Run("Success - 200 - CHILDREN", func(t *testing.T) {
		reports.EXPECT().ClassificationReport(gomock.Eq(enum.CHILDREN))
		runTest(app, buildGet("/reports/members/classification/children")).assertStatus(t, http.StatusOK)
	})
	t.Run("Success - 200 - TEEN", func(t *testing.T) {
		reports.EXPECT().ClassificationReport(gomock.Eq(enum.TEEN))
		runTest(app, buildGet("/reports/members/classification/teen")).assertStatus(t, http.StatusOK)
	})
	t.Run("Success - 200 - ADULT", func(t *testing.T) {
		reports.EXPECT().ClassificationReport(gomock.Eq(enum.ADULT))
		runTest(app, buildGet("/reports/members/classification/adult")).assertStatus(t, http.StatusOK)
	})
	t.Run("Success - 200 - YOUNG", func(t *testing.T) {
		reports.EXPECT().ClassificationReport(gomock.Eq(enum.YOUNG))
		runTest(app, buildGet("/reports/members/classification/young")).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 400", func(t *testing.T) {
		runTest(app, buildGet("/reports/members/classification/X")).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		reports.EXPECT().ClassificationReport(gomock.Eq(enum.YOUNG)).Return([]byte{}, genericError)
		runTest(app, buildGet("/reports/members/classification/young")).assertStatus(t, http.StatusInternalServerError)
	})

}