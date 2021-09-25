package security

import (
	"github.com/spf13/viper"
	"net/http"
	"testing"

	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	mock_security "github.com/BrunoDM2943/church-members-api/platform/security/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_security.NewMockUserRepository(ctrl)
	service := NewAuthService(repo)

	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindUser(userName).Return(buildUser("", generatePassword(password)), nil)
		token, err := service.GenerateToken(userName, password)
		assert.NotEmpty(t, token)
		assert.Nil(t, err)
	})
	t.Run("Fail - Not same password", func(t *testing.T) {
		repo.EXPECT().FindUser(userName).Return(buildUser("", generatePassword(password)), nil)
		token, err := service.GenerateToken(userName, password+"123")
		assert.Empty(t, token)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Error on Repo", func(t *testing.T) {
		repo.EXPECT().FindUser(userName).Return(nil, genericError)
		token, err := service.GenerateToken(userName, password)
		assert.Empty(t, token)
		assert.NotNil(t, err)
	})
	t.Run("Not found", func(t *testing.T) {
		repo.EXPECT().FindUser(userName).Return(nil, nil)
		token, err := service.GenerateToken(userName, password)
		assert.Empty(t, token)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.(apierrors.Error).StatusCode())
	})
}

func TestAuthService_IsValidToken(t *testing.T) {
	viper.Set("security.token.expiration", 1)
	assert.False(t, IsValidToken(""))
	assert.True(t, IsValidToken(buildToken(buildClaim())))
}

func generatePassword(password string) string {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(encrypted)
}
