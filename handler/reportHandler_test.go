package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_reports "github.com/BrunoDM2943/church-members-api/reports/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestGenerateBirthDayReportSucess(t *testing.T){
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reports := mock_reports.NewMockReportsGenerator(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(r)
	
	reports.EXPECT().BirthdayReport().Times(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reports/members/birthday", strings.NewReader(""))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestGenerateBirthDayReportFail(t *testing.T){
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reports := mock_reports.NewMockReportsGenerator(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(r)
	
	reports.EXPECT().BirthdayReport().Return(nil, errors.New(""))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reports/members/birthday", strings.NewReader(""))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fail()
	}
}

func TestGenerateMarriageReportSucess(t *testing.T){
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reports := mock_reports.NewMockReportsGenerator(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(r)
	
	reports.EXPECT().MariageReport().Times(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reports/members/marriage", strings.NewReader(""))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestGenerateMarriageReportFail(t *testing.T){
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reports := mock_reports.NewMockReportsGenerator(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(r)
	
	reports.EXPECT().MariageReport().Return(nil, errors.New(""))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reports/members/marriage", strings.NewReader(""))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fail()
	}
}