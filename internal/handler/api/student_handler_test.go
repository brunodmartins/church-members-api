package api

import (
	"net/http"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	mock_student "github.com/brunodmartins/church-members-api/internal/modules/student/mock"
	"go.uber.org/mock/gomock"
)

func TestPostStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := newApp()
	service := mock_student.NewMockService(ctrl)
	handler := NewStudentHandler(service)
	handler.SetUpRoutes(app)

	t.Run("creates a fresh student", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().CreateStudent(gomock.Any(), &domain.Student{Name: "Jane Doe", Gender: "F"}).Return(id, nil)

		runTest(app, buildPost("/students", []byte(`{"name":"Jane Doe","gender":"F"}`))).assertJSON(t, http.StatusCreated, dto.CreateMemberResponse{ID: id})
	})

	t.Run("creates a student from a member", func(t *testing.T) {
		id := domain.NewID()
		memberID := domain.NewID()
		service.EXPECT().CreateStudent(gomock.Any(), &domain.Student{Member: &domain.Member{ID: memberID}}).Return(id, nil)

		runTest(app, buildPost("/students", []byte(`{"memberId":"`+memberID+`"}`))).assertJSON(t, http.StatusCreated, dto.CreateMemberResponse{ID: id})
	})

	t.Run("creates a student from a participant", func(t *testing.T) {
		id := domain.NewID()
		participantID := domain.NewID()
		service.EXPECT().CreateStudent(gomock.Any(), &domain.Student{Participant: &domain.Participant{ID: participantID}}).Return(id, nil)

		runTest(app, buildPost("/students", []byte(`{"participantId":"`+participantID+`"}`))).assertJSON(t, http.StatusCreated, dto.CreateMemberResponse{ID: id})
	})

	t.Run("rejects an invalid source selection", func(t *testing.T) {
		runTest(app, buildPost("/students", []byte(`{}`))).assertStatus(t, http.StatusBadRequest)
		runTest(app, buildPost("/students", []byte(`{"memberId":"`+domain.NewID()+`","participantId":"`+domain.NewID()+`"}`))).assertStatus(t, http.StatusBadRequest)
	})
}
