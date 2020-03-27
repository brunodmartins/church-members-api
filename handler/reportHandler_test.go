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

func TestRoutesWithSuccess(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reports := mock_reports.NewMockReportsGenerator(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(r)

	reports.EXPECT().BirthdayReport().Times(1)
	reports.EXPECT().MariageReport().Times(1)
	reports.EXPECT().MemberReport().Times(1)

	routes := []string{
		"/reports/members/birthday",
		"/reports/members/marriage",
		"/reports/members",
	}

	for _, route := range routes {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", route, strings.NewReader(""))
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fail()
		}
	}
}

func TestRoutesWithFail(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reports := mock_reports.NewMockReportsGenerator(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(r)

	reports.EXPECT().BirthdayReport().Return(nil, errors.New("")).Times(1)
	reports.EXPECT().MariageReport().Return(nil, errors.New("")).Times(1)
	reports.EXPECT().MemberReport().Return(nil, errors.New("")).Times(1)

	routes := []string{
		"/reports/members/birthday",
		"/reports/members/marriage",
		"/reports/members",
	}

	for _, route := range routes {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", route, strings.NewReader(""))
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fail()
		}
	}
}
