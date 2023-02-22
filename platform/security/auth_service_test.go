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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_user.NewMockRepository(ctrl)
	churchRepo := mock_church.NewMockRepository(ctrl)
	service := NewAuthService(repo, churchRepo)

	t.Run("Success", func(t *testing.T) {
		churchRepo.EXPECT().GetByID(gomock.Any()).Return(&domain.Church{}, nil)
		repo.EXPECT().FindUser(gomock.Any(), userName).Return(buildUser("", string(crypto.EncryptPassword(password))), nil)
		token, err := service.GenerateToken("", userName, password)
		assert.NotEmpty(t, token)
		assert.Nil(t, err)
	})
	t.Run("Error on get church", func(t *testing.T) {
		churchRepo.EXPECT().GetByID(gomock.Any()).Return(&domain.Church{}, genericError)
		token, err := service.GenerateToken("", userName, password)
		assert.Empty(t, token)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Not same password", func(t *testing.T) {
		churchRepo.EXPECT().GetByID(gomock.Any()).Return(&domain.Church{}, nil)
		repo.EXPECT().FindUser(gomock.Any(), userName).Return(buildUser("", string(crypto.EncryptPassword(password))), nil)
		token, err := service.GenerateToken("", userName, password+"123")
		assert.Empty(t, token)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Error on Repo", func(t *testing.T) {
		churchRepo.EXPECT().GetByID(gomock.Any()).Return(&domain.Church{}, nil)
		repo.EXPECT().FindUser(gomock.Any(), userName).Return(nil, genericError)
		token, err := service.GenerateToken("", userName, password)
		assert.Empty(t, token)
		assert.NotNil(t, err)
	})
	t.Run("Not found", func(t *testing.T) {
		churchRepo.EXPECT().GetByID(gomock.Any()).Return(&domain.Church{}, nil)
		repo.EXPECT().FindUser(gomock.Any(), userName).Return(nil, nil)
		token, err := service.GenerateToken("", userName, password)
		assert.Empty(t, token)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.(apierrors.Error).StatusCode())
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
