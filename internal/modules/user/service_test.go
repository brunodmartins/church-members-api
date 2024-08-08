package user

import (
	"context"
	"errors"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"net/http"
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
	t.Run("Success", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Any(), gomock.Eq(user.UserName)).Return(nil, nil)
		repository.EXPECT().SaveUser(gomock.Any(), gomock.Eq(user)).Return(nil)
		assert.Nil(t, service.SaveUser(BuildContext(), user))
	})
	t.Run("Fail", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Any(), gomock.Eq(user.UserName)).Return(nil, nil)
		repository.EXPECT().SaveUser(gomock.Any(), gomock.Eq(user)).Return(genericError)
		assert.NotNil(t, service.SaveUser(BuildContext(), user))
	})
	t.Run("Fail - checking user - error", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Any(), gomock.Eq(user.UserName)).Return(nil, genericError)
		assert.NotNil(t, service.SaveUser(BuildContext(), user))
	})
	t.Run("Fail - checking user - already exist", func(t *testing.T) {
		repository.EXPECT().FindUser(gomock.Any(), gomock.Eq(user.UserName)).Return(user, nil)
		assert.NotNil(t, service.SaveUser(BuildContext(), user))
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

func TestUserService_ConfirmEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mock_user.NewMockRepository(ctrl)
	service := NewService(repository)
	t.Run("Given a valid user and token, when the confirm operation is call, then confirm correctly", func(t *testing.T) {
		id := domain.NewID()
		ctx := BuildContext()
		user := buildUser(id, "")
		user.ConfirmedEmail = false
		repository.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(user, nil)
		repository.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(user)).Return(nil)
		assert.Nil(t, service.ConfirmEmail(ctx, id, confirmationToken))
		assert.True(t, user.ConfirmedEmail)
	})
	t.Run("Given a valid user and token, when the confirm operation is call, then fails the operation due to repository error", func(t *testing.T) {
		id := domain.NewID()
		ctx := BuildContext()
		user := buildUser(id, "")
		user.ConfirmedEmail = false
		repository.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(user, nil)
		repository.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(user)).Return(errors.New("generic error"))
		assert.Error(t, service.ConfirmEmail(ctx, id, confirmationToken))
	})
	t.Run("Given a valid user and an invalid token, when the confirm operation is call, then fails the operation due to invalid token", func(t *testing.T) {
		id := domain.NewID()
		ctx := BuildContext()
		user := buildUser(id, "")
		user.ConfirmedEmail = false
		repository.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(user, nil)
		err := service.ConfirmEmail(ctx, id, "invalid-token")
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(apierrors.Error).StatusCode())
	})
	t.Run("Given a user that has already confirmed, when the confirm operation is call, then fails the operation due to already be confirmed", func(t *testing.T) {
		id := domain.NewID()
		ctx := BuildContext()
		user := buildUser(id, "")
		user.ConfirmedEmail = true
		repository.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(user, nil)
		err := service.ConfirmEmail(ctx, id, confirmationToken)
		assert.Error(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, err.(apierrors.Error).StatusCode())
	})
	t.Run("Given a nonexistent user, when the confirm operation is call, then fails the operation due to not found the user", func(t *testing.T) {
		id := domain.NewID()
		ctx := BuildContext()
		user := buildUser(id, "")
		user.ConfirmedEmail = true
		repository.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(nil, apierrors.NewApiError("User not found", http.StatusNotFound))
		err := service.ConfirmEmail(ctx, id, confirmationToken)
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, err.(apierrors.Error).StatusCode())
	})

}

func BuildContext() context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{
			ID: "church_id_test",
		},
	})
}
