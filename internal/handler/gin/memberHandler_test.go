package gin

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BrunoDM2943/church-members-api/internal/repository"

	"github.com/BrunoDM2943/church-members-api/internal/constants/entity"
	mock_service "github.com/BrunoDM2943/church-members-api/internal/service/member/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestGetMemberBadRequest(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)

	memberHandler.SetUpRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members/a", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestGetMemberNotFound(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)

	id := entity.NewID()

	service.EXPECT().FindMembersByID(id).Return(nil, repository.MemberNotFound)
	memberHandler.SetUpRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members/"+id, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fail()
	}
}

func TestGetMemberOK(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)

	member := &entity.Member{}
	member.ID = entity.NewID()
	service.EXPECT().FindMembersByID(member.ID).Return(member, nil).AnyTimes()

	memberHandler.SetUpRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members/"+member.ID, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestPostMemberBadRequest(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/members", strings.NewReader(""))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestPostMemberSucess(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	body := `{"attendsFridayWorship":true,"attendsSaturdayWorship":true,"attendsSundayWorship":true,"attendsSundaySchool":true,"person":{"firstName":"XXXX","lastName":"XXXX XXX","birthDate":"2020-01-01T00:00:00-03:00","naturalidade":"XXXXX","placeOfBirth":"XXXXX","fathersName":"XXXXXX","mothersName":"XXXXX","spousesName":"XXXXX","brothersQuantity":0,"childrensQuantity":0,"profession":"XXXXXX","gender":"M","contact":{"cellPhoneArea":99,"cellPhone":123456789,"email":"XXXXXX"},"address":{"zipCode":"XXXXXXX","state":"XXXX","city":"XXXXX","address":"XXXX","district":"XXX","number":1,"moreInfo":"xxXXX"}},"religion":{"fathersReligion":"Crentes","baptismPlace":"IEPEM","learnedGospelAge":10,"acceptedJesus":true,"baptized":true,"catholicBaptized":true,"knowsTithe":true,"agreesTithe":true,"tithe":true}}`
	service.EXPECT().SaveMember(gomock.AssignableToTypeOf(&entity.Member{})).Return(entity.NewID(), nil)
	req, _ := http.NewRequest("POST", "/members", strings.NewReader(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		log.Print(string(w.Body.Bytes()))
		t.Fail()
	}
}

func TestPostMemberFail(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	body := `{"attendsFridayWorship":true,"attendsSaturdayWorship":true,"attendsSundayWorship":true,"attendsSundaySchool":true,"person":{"firstName":"XXXX","lastName":"XXXX XXX","birthDate":"2020-01-01T00:00:00-03:00","naturalidade":"XXXXX","placeOfBirth":"XXXXX","fathersName":"XXXXXX","mothersName":"XXXXX","spousesName":"XXXXX","brothersQuantity":0,"childrensQuantity":0,"profession":"XXXXXX","gender":"M","contact":{"cellPhoneArea":99,"cellPhone":123456789,"email":"XXXXXX"},"address":{"zipCode":"XXXXXXX","state":"XXXX","city":"XXXXX","address":"XXXX","district":"XXX","number":1,"moreInfo":"xxXXX"}},"religion":{"fathersReligion":"Crentes","baptismPlace":"IEPEM","learnedGospelAge":10,"acceptedJesus":true,"baptized":true,"catholicBaptized":true,"knowsTithe":true,"agreesTithe":true,"tithe":true}}`
	service.EXPECT().SaveMember(gomock.Any()).Return(entity.NewID(), errors.New(""))
	req, _ := http.NewRequest("POST", "/members", strings.NewReader(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fail()
	}
}

func TestPostMemberSearch(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	body := `
	{
		member(gender:"M", active:false){
				person{
					firstName,
					lastName,
					gender
				}
		}
	}`
	service.EXPECT().FindMembers(gomock.Any()).Return([]*entity.Member{}, nil)
	req, _ := http.NewRequest("POST", "/members/search", strings.NewReader(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestPostMemberSearchError(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	body := `
	{
		member(gender:"M", active:false){
				person{
					name,
					lastNa
					gender
				}
		}
	}`
	req, _ := http.NewRequest("POST", "/members/search", strings.NewReader(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fail()
	}
}

func TestPutStatus(t *testing.T) {
	type data struct {
		url        string
		body       string
		statusCode int
	}
	id := entity.NewID()
	urlWithID := "/members/" + id + "/status"
	table := []data{
		{"/members/X/status", "", http.StatusBadRequest},
		{urlWithID, `{"active":false}`, http.StatusBadRequest},
		{urlWithID, `{"reason": "exited"}`, http.StatusBadRequest},
		{urlWithID, `{"active":false, "reason": "exited"}`, http.StatusInternalServerError},
		{urlWithID, `{"active":true, "reason": "Comed back"}`, http.StatusOK}}

	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	service.EXPECT().ChangeStatus(id, gomock.Eq(false), gomock.Eq("exited"), gomock.Any()).Return(errors.New("Error"))
	service.EXPECT().ChangeStatus(id, gomock.Eq(true), gomock.Eq("Comed back"), gomock.Any()).Return(nil)

	memberHandler.SetUpRoutes(r)

	for _, test := range table {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", test.url, strings.NewReader(test.body))
		r.ServeHTTP(w, req)
		if w.Code != test.statusCode {
			t.Errorf("Failed for test: %s, %s, Status Code: %d, Expected: %d", test.url, test.body, w.Code, test.statusCode)
		}
	}
}
