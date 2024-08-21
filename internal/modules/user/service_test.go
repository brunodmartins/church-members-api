package user

import (
	"context"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	mock_user "github.com/brunodmartins/church-members-api/internal/modules/user/mock"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_SaveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mock_user.NewMockRepository(ctrl)
	service := NewService(repository)
	user := buildUser("id", "")
	t.Run("Given a valid user, when save it, then store on the database successfully", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Any(), gomock.Eq(user.UserName)).Return(nil, nil)
		repository.EXPECT().SaveUser(gomock.Any(), gomock.Eq(user)).Return(nil)
		assert.NoError(t, service.SaveUser(BuildContext(), user))
	})
	t.Run("Given a valid user, when save it, then store on the database fails", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Any(), gomock.Eq(user.UserName)).Return(nil, nil)
		repository.EXPECT().SaveUser(gomock.Any(), gomock.Eq(user)).Return(genericError)
		assert.Error(t, service.SaveUser(BuildContext(), user))
	})
	t.Run("Given a valid user, when save it, then fails upon validating duplication", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Any(), gomock.Eq(user.UserName)).Return(nil, genericError)
		assert.Error(t, service.SaveUser(BuildContext(), user))
	})
	t.Run("Given a valid user, when save it, then fails due to already exist", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Any(), gomock.Eq(user.UserName)).Return(user, nil)
		assert.Error(t, service.SaveUser(BuildContext(), user))
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
		repository.EXPECT().SearchUser(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return([]*domain.User{user}, nil)
		result, err := service.SearchUser(BuildContext(), spec)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})
	t.Run("Fail", func(t *testing.T) {
		repository.EXPECT().SearchUser(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return([]*domain.User{}, genericError)
		_, err := service.SearchUser(BuildContext(), spec)
		assert.NotNil(t, err)
	})
}

func TestService_FindUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mock_user.NewMockRepository(ctrl)
	service := NewService(repository)
	user := buildUser("id", "")
	ctx := BuildContext()
	t.Run("Given a valid username, when performing a find on the database, then return the user successfully", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Eq(ctx), gomock.Eq(userName)).Return(user, nil)
		result, err := service.FindUser(ctx, userName)
		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})
	t.Run("Given a valid username, when performing a find on the database, then return error ", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Eq(ctx), gomock.Eq(userName)).Return(nil, genericError)
		_, err := service.FindUser(ctx, userName)
		assert.Error(t, err)
	})
}

func TestService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mock_user.NewMockRepository(ctrl)
	service := NewService(repository)
	user := buildUser("id", "")
	ctx := BuildContext()
	t.Run("Given a valid user, when performing an update on the database, then return no error", func(t *testing.T) {
		repository.EXPECT().UpdateUser(gomock.Eq(ctx), gomock.Eq(user)).Return(nil)
		assert.NoError(t, service.UpdateUser(ctx, user))
	})
	t.Run("Given a valid username, when performing an update on the database, then return error ", func(t *testing.T) {
		repository.EXPECT().UpdateUser(gomock.Eq(ctx), gomock.Eq(user)).Return(genericError)
		assert.Error(t, service.UpdateUser(ctx, user))
	})
}

func BuildContext() context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{
			ID: "church_id_test",
		},
	})
}
