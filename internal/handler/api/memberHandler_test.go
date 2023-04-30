package api

import (
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	mock_member "github.com/brunodmartins/church-members-api/internal/modules/member/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
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
		service.EXPECT().GetMember(gomock.Any(), id).Return(buildMember(id), nil)
		runTest(app, buildGet("/members/"+id)).assert(t, http.StatusOK, new(dto.GetMemberResponse), func(parsedBody interface{}) {
			member := parsedBody.(*dto.GetMemberResponse)
			assert.Equal(t, id, member.ID)
		})
	})
	t.Run("Fail - 404", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetMember(gomock.Any(), id).Return(nil, apierrors.NewApiError("Member not found", http.StatusNotFound))
		runTest(app, buildGet("/members/"+id)).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 400", func(t *testing.T) {
		runTest(app, buildGet("/members/a")).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetMember(gomock.Any(), id).Return(nil, genericError)
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
		service.EXPECT().SaveMember(gomock.Any(), gomock.AssignableToTypeOf(&domain.Member{})).Return(id, nil)
		runTest(app, buildPost("/members", body)).assert(t, http.StatusCreated, new(dto.CreateMemberResponse), func(parsedBody interface{}) {
			response := parsedBody.(*dto.CreateMemberResponse)
			assert.Equal(t, id, response.ID)
		})
	})
	t.Run("Fail - 404", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetMember(gomock.Any(), id).Return(nil, apierrors.NewApiError("Member not found", http.StatusNotFound))
		runTest(app, buildGet("/members/"+id)).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 400", func(t *testing.T) {
		runTest(app, buildPost("/members", emptyJson)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		id := domain.NewID()
		body, _ := os.ReadFile("./resources/create_member.json")
		service.EXPECT().SaveMember(gomock.Any(), gomock.AssignableToTypeOf(&domain.Member{})).Return(id, genericError)
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
		service.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).Return([]*domain.Member{}, nil)
		runTest(app, buildGet("/members?name=test&active=true")).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		service.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).Return([]*domain.Member{}, genericError)
		runTest(app, buildGet("/members?name=test&active=true")).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestRetireMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_member.NewMockService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(app)
	id := domain.NewID()

	t.Run("Success - 200", func(t *testing.T) {
		body := []byte(`{"reason": "Leaved the church"}`)
		service.EXPECT().RetireMembership(gomock.Any(), id, gomock.Eq("Leaved the church"), gomock.Any()).Return(nil)
		runTest(app, buildDelete(fmt.Sprintf("/members/%s", id), body)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 400 - ID", func(t *testing.T) {
		runTest(app, buildDelete(fmt.Sprintf("/members/%s", "X"), emptyJson)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 400 - Reason", func(t *testing.T) {
		body := []byte(`{}`)
		runTest(app, buildDelete(fmt.Sprintf("/members/%s", id), body)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 404", func(t *testing.T) {
		body := []byte(`{"reason": "Not Found"}`)
		service.EXPECT().RetireMembership(gomock.Any(), id, gomock.Eq("Not Found"), gomock.Any()).Return(apierrors.NewApiError("Member not found", http.StatusNotFound))
		runTest(app, buildDelete(fmt.Sprintf("/members/%s", id), body)).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		body := []byte(`{"reason": "exited"}`)
		service.EXPECT().RetireMembership(gomock.Any(), id, gomock.Eq("exited"), gomock.Any()).Return(genericError)
		runTest(app, buildDelete(fmt.Sprintf("/members/%s", id), body)).assertStatus(t, http.StatusInternalServerError)
	})
}
