package church

import (
	"errors"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	mock_church "github.com/brunodmartins/church-members-api/internal/modules/church/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

func TestChurchService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_church.NewMockRepository(ctrl)
	service := NewService(repo)

	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().List().Return([]*domain.Church{buildChurch("")}, nil)
		result, err := service.List()
		assert.Nil(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().List().Return(nil, genericError)
		_, err := service.List()
		assert.NotNil(t, err)
	})
}

func TestChurchService_GetChurch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_church.NewMockRepository(ctrl)
	service := NewService(repo)
	const id = "xxx"
	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().GetByID(id).Return(buildChurch(id), nil)
		result, err := service.GetChurch(id)
		assert.Nil(t, err)
		assert.Equal(t, id, result.ID)
	})

	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().GetByID(id).Return(nil, genericError)
		_, err := service.GetChurch(id)
		assert.NotNil(t, err)
	})
}

func TestChurchService_GetChurchByAbbreviation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_church.NewMockRepository(ctrl)
	service := NewService(repo)
	const abbreviation = "TC"
	const id = "xxx"

	t.Run("Success - Church Found", func(t *testing.T) {
		expected := buildChurch(id)
		expected.Abbreviation = abbreviation
		churches := []*domain.Church{expected}

		repo.EXPECT().List().Return(churches, nil)

		result, err := service.GetChurchByAbbreviation(abbreviation)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		repo.EXPECT().List().Return(nil, errors.New("database error"))

		result, err := service.GetChurchByAbbreviation(abbreviation)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error - Church Not Found", func(t *testing.T) {
		expected := buildChurch(id)
		expected.Abbreviation = abbreviation
		churches := []*domain.Church{expected}

		repo.EXPECT().List().Return(churches, nil)

		result, err := service.GetChurchByAbbreviation("OTHER")
		assert.Error(t, err)
		assert.Nil(t, result)

		// Check if it's an API error with correct status
		apiErr, ok := err.(apierrors.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})
}
