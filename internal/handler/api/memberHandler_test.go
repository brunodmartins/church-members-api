package api

import (
	"fmt"
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member/mock"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/golang/mock/gomock"
)

func TestGetMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_member.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetMember(id).Return(buildMember(id), nil)
		runTest(app, buildGet("/members/"+id)).assert(t, http.StatusOK, new(domain.Member), func(parsedBody interface{}) {
			member := parsedBody.(*domain.Member)
			assert.Equal(t, id, member.ID)
		})
	})
	t.Run("Fail - 404", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetMember(id).Return(nil, apierrors.NewApiError("Member not found", http.StatusNotFound))
		runTest(app, buildGet("/members/"+id)).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 400", func(t *testing.T) {
		runTest(app, buildGet("/members/a")).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetMember(id).Return(nil, genericError)
		runTest(app, buildGet("/members/"+id)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestPostMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_member.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		id := domain.NewID()
		body := getMock("create_member.json")
		service.EXPECT().SaveMember(gomock.AssignableToTypeOf(&domain.Member{})).Return(id, nil)
		runTest(app, buildPost("/members", body)).assert(t, http.StatusCreated, new(dto.CreateMemberResponse), func(parsedBody interface{}) {
			response := parsedBody.(*dto.CreateMemberResponse)
			assert.Equal(t, id, response.ID)
		})
	})
	t.Run("Fail - 404", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetMember(id).Return(nil, apierrors.NewApiError("Member not found", http.StatusNotFound))
		runTest(app, buildGet("/members/"+id)).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 400", func(t *testing.T) {
		runTest(app, buildPost("/members", emptyJson)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		id := domain.NewID()
		body, _ := ioutil.ReadFile("./resources/create_member.json")
		service.EXPECT().SaveMember(gomock.AssignableToTypeOf(&domain.Member{})).Return(id, genericError)
		runTest(app, buildPost("/members", body)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestPostMemberSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_member.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		body := []byte(`{
			member(gender:"M", active:false){
					person{
						firstName,
						lastName,
						gender
					}
			}
		}`)
		service.EXPECT().SearchMembers(gomock.Any()).Return([]*domain.Member{}, nil)
		runTest(app, buildPost("/members/search", body)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		runTest(app, buildPost("/members/search", emptyJson)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestPutStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_member.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(app)
	id := domain.NewID()

	t.Run("Success - 200", func(t *testing.T) {
		body := []byte(`{"active":true, "reason": "Came back"}`)
		service.EXPECT().ChangeStatus(id, gomock.Eq(true), gomock.Eq("Came back"), gomock.Any()).Return(nil)
		runTest(app, buildPut(fmt.Sprintf("/members/%s/status", id), body)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 400 - ID", func(t *testing.T) {
		runTest(app, buildPut(fmt.Sprintf("/members/%s/status", "X"), emptyJson)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 400 - Reason", func(t *testing.T) {
		body := []byte(`{"active":false}`)
		runTest(app, buildPut(fmt.Sprintf("/members/%s/status", id), body)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 400 - Active", func(t *testing.T) {
		body := []byte(`{"reason": "exited"}`)
		runTest(app, buildPut(fmt.Sprintf("/members/%s/status", id), body)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 404", func(t *testing.T) {
		body := []byte(`{"active":false, "reason": "Not Found"}`)
		service.EXPECT().ChangeStatus(id, gomock.Eq(false), gomock.Eq("Not Found"), gomock.Any()).Return(apierrors.NewApiError("Member not found", http.StatusNotFound))
		runTest(app, buildPut(fmt.Sprintf("/members/%s/status", id), body)).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		body := []byte(`{"active":false, "reason": "exited"}`)
		service.EXPECT().ChangeStatus(id, gomock.Eq(false), gomock.Eq("exited"), gomock.Any()).Return(genericError)
		runTest(app, buildPut(fmt.Sprintf("/members/%s/status", id), body)).assertStatus(t, http.StatusInternalServerError)
	})
}
