package church

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	mock_church "github.com/BrunoDM2943/church-members-api/internal/modules/church/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
