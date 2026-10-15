package confirmation

import (
	"testing"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	viper.Set("security.token.secret", "test-confirmation-secret")
	user := &domain.User{
		ID:                "user-id",
		ChurchID:          "church-id",
		ConfirmationToken: "confirmation-nonce",
	}

	before := time.Now()
	token, err := GenerateToken(user)
	require.NoError(t, err)

	claims := new(confirmationClaim)
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})
	require.NoError(t, err)
	assert.True(t, parsed.Valid)
	assert.WithinDuration(t, before.Add(time.Hour), claims.ExpiresAt.Time, time.Second)
	assert.NoError(t, ValidateToken(token, user))

	otherUser := *user
	otherUser.ID = "another-user"
	assert.Error(t, ValidateToken(token, &otherUser))
	assert.Error(t, ValidateToken(token+"tampered", user))
}

func TestTokenRequiresSigningSecret(t *testing.T) {
	viper.Set("security.token.secret", "")
	_, err := GenerateToken(&domain.User{ID: "user-id"})
	assert.Error(t, err)
	assert.Error(t, ValidateToken("token", &domain.User{ID: "user-id"}))
}

func TestValidateTokenRejectsExpiredToken(t *testing.T) {
	viper.Set("security.token.secret", "test-confirmation-secret")
	user := &domain.User{
		ID:                "user-id",
		ChurchID:          "church-id",
		ConfirmationToken: "confirmation-nonce",
	}
	claims := confirmationClaim{
		ChurchID:          user.ChurchID,
		ConfirmationToken: user.ConfirmationToken,
		Purpose:           confirmationTokenPurpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Second)),
			Issuer:    "church-members-api",
			Subject:   user.ID,
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(getSecret())
	require.NoError(t, err)

	assert.Error(t, ValidateToken(token, user))
}
