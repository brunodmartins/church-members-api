package api

import (
	"net/http"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	mock_church "github.com/brunodmartins/church-members-api/internal/modules/church/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetChurch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_church.NewMockService(ctrl)
	handler := NewChurchHandler(service)
	handler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetChurch(gomock.Any(), id).Return(buildChurch(id), nil)
		runTest(app, buildGet("/churches/"+id)).assert(t, http.StatusOK, new(dto.GetChurchResponse), func(parsedBody interface{}) {
			response := parsedBody.(*dto.GetChurchResponse)
			assert.Equal(t, id, response.ID)
		})
	})
	t.Run("Fail - 404", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetChurch(gomock.Any(), id).Return(nil, apierrors.NewApiError("Member not found", http.StatusNotFound))
		runTest(app, buildGet("/churches/"+id)).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetChurch(gomock.Any(), id).Return(nil, genericError)
		runTest(app, buildGet("/churches/"+id)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestGetStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_church.NewMockService(ctrl)
	handler := NewChurchHandler(service)
	handler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetStatistics(gomock.Any(), id).Return(buildStatistics(), nil)
		runTest(app, buildGet("/churches/"+id+"/statistics")).assert(t, http.StatusOK, new(dto.ChurchStatisticsResponse), func(parsedBody interface{}) {
			response := parsedBody.(*dto.ChurchStatisticsResponse)
			assert.Equal(t, 10, response.TotalMembers)
		})
	})
	t.Run("Fail - 404", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetStatistics(gomock.Any(), id).Return(nil, apierrors.NewApiError("Church not found", http.StatusNotFound))
		runTest(app, buildGet("/churches/"+id+"/statistics")).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetStatistics(gomock.Any(), id).Return(nil, genericError)
		runTest(app, buildGet("/churches/"+id+"/statistics")).assertStatus(t, http.StatusInternalServerError)
	})
}

func buildStatistics() *domain.ChurchStatistics {
	return &domain.ChurchStatistics{
		TotalMembers:                 10,
		AgeDistribution:              []int{1, 2, 3},
		TotalMembersByGender:         map[string]int{"M": 5, "F": 5},
		TotalMembersByClassification: map[string]int{"MEMBER": 10},
	}
}
