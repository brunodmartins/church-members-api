package user

import (
	mock_user "github.com/BrunoDM2943/church-members-api/internal/modules/user/mock"
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
		assert.Nil(t, service.SaveUser(user))
	})
	t.Run("Fail", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Eq(user.UserName)).Return(nil, nil)
		repository.EXPECT().SaveUser(gomock.Eq(user)).Return(genericError)
		assert.NotNil(t, service.SaveUser(user))
	})
	t.Run("Fail - checking user - error", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Eq(user.UserName)).Return(nil, genericError)
		assert.NotNil(t, service.SaveUser(user))
	})
	t.Run("Fail - checking user - already exist", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Eq(user.UserName)).Return(user, nil)
		assert.NotNil(t, service.SaveUser(user))
	})
}
