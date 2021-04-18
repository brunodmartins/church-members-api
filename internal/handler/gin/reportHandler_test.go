package gin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_report "github.com/BrunoDM2943/church-members-api/internal/service/report/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestRoutesWithSuccess(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(r)

	reports.EXPECT().BirthdayReport().Times(1)
	reports.EXPECT().MarriageReport().Times(1)
	reports.EXPECT().MemberReport().Times(1)
	reports.EXPECT().LegalReport().Times(1)
	reports.EXPECT().ClassificationReport(gomock.Eq("children")).Times(1)
	reports.EXPECT().ClassificationReport(gomock.Eq("teen")).Times(1)
	reports.EXPECT().ClassificationReport(gomock.Eq("adult")).Times(1)
	reports.EXPECT().ClassificationReport(gomock.Eq("young")).Times(1)

	routes := []string{
		"/reports/members/birthday",
		"/reports/members/marriage",
		"/reports/members/legal",
		"/reports/members/classification/children",
		"/reports/members/classification/teen",
		"/reports/members/classification/adult",
		"/reports/members/classification/young",
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

	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(r)

	reports.EXPECT().BirthdayReport().Return(nil, errors.New("")).Times(1)
	reports.EXPECT().MarriageReport().Return(nil, errors.New("")).Times(1)
	reports.EXPECT().MemberReport().Return(nil, errors.New("")).Times(1)
	reports.EXPECT().ClassificationReport("adult").Return(nil, errors.New("")).Times(1)

	routes := []string{
		"/reports/members/birthday",
		"/reports/members/marriage",
		"/reports/members/classification/adult",
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

func TestRoutesForClassificationWithBadRequest(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reports := mock_report.NewMockService(ctrl)
	reportHandler := NewReportHandler(reports)
	reportHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reports/members/classification/X", strings.NewReader(""))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}
