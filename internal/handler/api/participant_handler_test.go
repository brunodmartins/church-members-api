package api

import (
	"net/http"
	"os"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	mock_participant "github.com/brunodmartins/church-members-api/internal/modules/participant/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

func buildParticipant(id string) *domain.Participant {
	return &domain.Participant{
		ID:     id,
		Name:   "Test Participant",
		Gender: "MALE",
	}
}

func TestGetParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_participant.NewMockService(ctrl)
	participantHandler := NewParticipantHandler(service)
	participantHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetParticipant(gomock.Any(), id).Return(buildParticipant(id), nil)
		runTest(app, buildGet("/participants/"+id)).assert(t, http.StatusOK, new(dto.GetParticipantResponse), func(parsedBody interface{}) {
			p := parsedBody.(*dto.GetParticipantResponse)
			assert.Equal(t, id, p.ID)
		})
	})
	t.Run("Fail - 404", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetParticipant(gomock.Any(), id).Return(nil, apierrors.NewApiError("Participant not found", http.StatusNotFound))
		runTest(app, buildGet("/participants/"+id)).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 400", func(t *testing.T) {
		runTest(app, buildGet("/participants/a")).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetParticipant(gomock.Any(), id).Return(nil, genericError)
		runTest(app, buildGet("/participants/"+id)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestPostParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_participant.NewMockService(ctrl)
	participantHandler := NewParticipantHandler(service)
	participantHandler.SetUpRoutes(app)

	t.Run("Success - 201", func(t *testing.T) {
		id := domain.NewID()
		body := getMock("create_participant.json")
		service.EXPECT().CreateParticipant(gomock.Any(), gomock.AssignableToTypeOf(&domain.Participant{})).Return(id, nil)
		runTest(app, buildPost("/participants", body)).assert(t, http.StatusCreated, new(dto.CreateMemberResponse), func(parsedBody interface{}) {
			response := parsedBody.(*dto.CreateMemberResponse)
			assert.Equal(t, id, response.ID)
		})
	})
	t.Run("Fail - 400", func(t *testing.T) {
		runTest(app, buildPost("/participants", emptyJson)).assertStatus(t, http.StatusBadRequest)
		runTest(app, buildPost("/participants", badJson)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		id := domain.NewID()
		body, _ := os.ReadFile("./resources/create_participant.json")
		service.EXPECT().CreateParticipant(gomock.Any(), gomock.AssignableToTypeOf(&domain.Participant{})).Return(id, genericError)
		runTest(app, buildPost("/participants", body)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestSearchParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_participant.NewMockService(ctrl)
	participantHandler := NewParticipantHandler(service)
	participantHandler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		service.EXPECT().SearchParticipant(gomock.Any(), gomock.Any()).Return([]*domain.Participant{}, nil)
		runTest(app, buildGet("/participants?name=test&gender=MALE")).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		service.EXPECT().SearchParticipant(gomock.Any(), gomock.Any()).Return([]*domain.Participant{}, genericError)
		runTest(app, buildGet("/participants?name=test&gender=MALE")).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestUpdateParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_participant.NewMockService(ctrl)
	participantHandler := NewParticipantHandler(service)
	participantHandler.SetUpRoutes(app)
	id := domain.NewID()

	t.Run("Success - 200", func(t *testing.T) {
		body := getMock("create_participant.json")
		service.EXPECT().UpdateParticipant(gomock.Any(), gomock.AssignableToTypeOf(&domain.Participant{})).Return(nil)
		runTest(app, buildPut("/participants/"+id, body)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 400 - ID", func(t *testing.T) {
		runTest(app, buildPut("/participants/X", emptyJson)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 400 - Bad JSON", func(t *testing.T) {
		runTest(app, buildPut("/participants/"+id, badJson)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		body := getMock("create_participant.json")
		service.EXPECT().UpdateParticipant(gomock.Any(), gomock.AssignableToTypeOf(&domain.Participant{})).Return(genericError)
		runTest(app, buildPut("/participants/"+id, body)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestDeleteParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_participant.NewMockService(ctrl)
	participantHandler := NewParticipantHandler(service)
	participantHandler.SetUpRoutes(app)
	id := domain.NewID()

	t.Run("Success - 200", func(t *testing.T) {
		service.EXPECT().DeleteParticipant(gomock.Any(), id).Return(nil)
		runTest(app, buildDelete("/participants/"+id, nil)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - 400 - ID", func(t *testing.T) {
		runTest(app, buildDelete("/participants/X", nil)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - 404", func(t *testing.T) {
		service.EXPECT().DeleteParticipant(gomock.Any(), id).Return(apierrors.NewApiError("Participant not found", http.StatusNotFound))
		runTest(app, buildDelete("/participants/"+id, nil)).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		service.EXPECT().DeleteParticipant(gomock.Any(), id).Return(genericError)
		runTest(app, buildDelete("/participants/"+id, nil)).assertStatus(t, http.StatusInternalServerError)
	})
}
