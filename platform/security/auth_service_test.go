package security

import (
	"context"
	"net/http"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	mock_church "github.com/brunodmartins/church-members-api/internal/modules/church/mock"
	mock_user "github.com/brunodmartins/church-members-api/internal/modules/user/mock"
	"github.com/brunodmartins/church-members-api/platform/crypto"
	"github.com/spf13/viper"

	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mock_user.NewMockService(ctrl)
	churchService := mock_church.NewMockService(ctrl)
	service := NewAuthService(userService, churchService)
	church := buildChurch(domain.NewID())
	testUser := buildUser("", string(crypto.EncryptPassword(password)))
	t.Run("Given a valid user, when try to authenticate, then succeeds", func(t *testing.T) {
		churchService.EXPECT().GetChurch(gomock.Eq(church.ID)).Return(church, nil)
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return([]*domain.User{
			testUser,
		}, nil)
		token, err := service.GenerateToken(church.ID, userName, password)
		assert.NotEmpty(t, token)
		assert.Nil(t, err)
	})
	t.Run("Given a valid user, when try to authenticate, then fails due to error on look up for the church", func(t *testing.T) {
		churchService.EXPECT().GetChurch(gomock.Eq(church.ID)).Return(nil, genericError)
		token, err := service.GenerateToken(church.ID, userName, password)
		assert.Empty(t, token)
		assert.Error(t, err)
	})
	t.Run("Given a valid user, when try to authenticate, then fails due to password not match", func(t *testing.T) {
		churchService.EXPECT().GetChurch(gomock.Eq(church.ID)).Return(church, nil)
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return([]*domain.User{
			testUser,
		}, nil)
		token, err := service.GenerateToken(church.ID, userName, password+"123")
		assert.Empty(t, token)
		assert.Error(t, err)
	})
	t.Run("Given a valid user, when try to authenticate, then fails due to error on look up for the user", func(t *testing.T) {
		churchService.EXPECT().GetChurch(gomock.Eq(church.ID)).Return(church, nil)
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(nil, genericError)
		token, err := service.GenerateToken(church.ID, userName, password)
		assert.Empty(t, token)
		assert.Error(t, err)
	})
	t.Run("Given a valid user, when try to authenticate, then fails due to user not been found", func(t *testing.T) {
		churchService.EXPECT().GetChurch(gomock.Eq(church.ID)).Return(church, nil)
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return([]*domain.User{}, nil)
		token, err := service.GenerateToken(church.ID, userName, password)
		assert.Empty(t, token)
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, err.(apierrors.Error).StatusCode())
	})
	t.Run("Given a valid user, when try to authenticate, then fails due to user not have confirmed the email, and send email", func(t *testing.T) {
		churchService.EXPECT().GetChurch(gomock.Eq(church.ID)).Return(church, nil)
		unconfirmedUser := buildUser("", string(crypto.EncryptPassword(password)))
		unconfirmedUser.ConfirmedEmail = false
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return([]*domain.User{unconfirmedUser}, nil)
		userService.EXPECT().SendConfirmEmail(gomock.Any(), gomock.Eq(unconfirmedUser))
		token, err := service.GenerateToken(church.ID, userName, password)
		assert.Empty(t, token)
		assert.Error(t, err)
	})
}

func TestAuthService_IsValidToken(t *testing.T) {
	viper.Set("security.token.expiration", 1)
	valid, _ := GetClaim("")
	assert.False(t, valid)
	valid, _ = GetClaim(buildToken())
	assert.True(t, valid)
}

func TestAddClaimToContext(t *testing.T) {
	viper.Set("security.token.expiration", 1)
	ctx := context.Background()
	_, claim := GetClaim(buildToken())
	assert.Nil(t, ctx.Value("user"))
	ctx = AddClaimToContext(claim, ctx)
	assert.NotNil(t, ctx.Value("user"))
}

func buildChurch(id string) *domain.Church {
	return &domain.Church{
		ID:           id,
		Name:         "test church",
		Abbreviation: "tc",
		Language:     "pt-br",
		Email:        "test@test.com",
	}
}
