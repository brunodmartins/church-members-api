package security

import (
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	mock_security "github.com/BrunoDM2943/church-members-api/platform/security/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAuthService_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_security.NewMockUserRepository(ctrl)
	service := NewAuthService(repo)

	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindUser(userName, gomock.Any()).Return(buildUser(""), nil)
		token, err := service.GenerateToken(userName, password)
		assert.NotEmpty(t, token)
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().FindUser(userName, gomock.Any()).Return(buildUser(""), genericError)
		token, err := service.GenerateToken(userName, password)
		assert.Empty(t, token)
		assert.NotNil(t, err)
	})
	t.Run("Not found", func(t *testing.T) {
		repo.EXPECT().FindUser(userName, gomock.Any()).Return(nil, nil)
		token, err := service.GenerateToken(userName, password)
		assert.Empty(t, token)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.(apierrors.Error).StatusCode())
	})
}


func TestAuthService_IsValidToken(t *testing.T) {
	assert.False(t, IsValidToken(""))
	assert.True(t, IsValidToken(buildToken(buildClaim())))
}