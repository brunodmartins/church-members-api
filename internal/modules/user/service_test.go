package user

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	mock_user "github.com/BrunoDM2943/church-members-api/internal/modules/user/mock"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_SaveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mock_user.NewMockRepository(ctrl)
	service := NewService(repository)
	user := buildUser("id", "")
	t.Run("Success", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Eq(user.UserName)).Return(nil, nil)
		repository.EXPECT().SaveUser(gomock.Eq(user)).Return(nil)
		assert.Nil(t, service.SaveUser(context.Background(), user))
	})
	t.Run("Fail", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Eq(user.UserName)).Return(nil, nil)
		repository.EXPECT().SaveUser(gomock.Eq(user)).Return(genericError)
		assert.NotNil(t, service.SaveUser(context.Background(), user))
	})
	t.Run("Fail - checking user - error", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Eq(user.UserName)).Return(nil, genericError)
		assert.NotNil(t, service.SaveUser(context.Background(), user))
	})
	t.Run("Fail - checking user - already exist", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Eq(user.UserName)).Return(user, nil)
		assert.NotNil(t, service.SaveUser(context.Background(), user))
	})
}

func TestService_SearchUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mock_user.NewMockRepository(ctrl)
	service := NewService(repository)
	user := buildUser("id", "")
	spec := wrapper.QuerySpecification(nil)
	t.Run("Success", func(t *testing.T) {
		repository.EXPECT().SearchUser(gomock.AssignableToTypeOf(spec)).Return([]*domain.User{user}, nil)
		result, err := service.SearchUser(context.Background(), spec)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})
	t.Run("Fail", func(t *testing.T) {
		repository.EXPECT().SearchUser(gomock.AssignableToTypeOf(spec)).Return([]*domain.User{}, genericError)
		_, err := service.SearchUser(context.Background(), spec)
		assert.NotNil(t, err)
	})
}
