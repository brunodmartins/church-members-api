package security

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_generateToken(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		jwtToken := GenerateJWTToken(buildClaim())
		assert.NotEmpty(t, jwtToken)
	})
}

func Test_getClaim(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		viper.Set("security.token.expiration", 1)
		claim, err := getClaim(buildToken(buildClaim()))
		assert.NotNil(t, claim)
		assert.Nil(t, err)
	})
	t.Run("Fail - empty", func(t *testing.T) {
		viper.Set("security.token.expiration", 1)
		claim, err := getClaim("")
		assert.Nil(t, claim)
		assert.NotNil(t, err)
	})
	t.Run("Fail - expired", func(t *testing.T) {
		viper.Set("security.token.expiration", -1)
		claim, err := getClaim(buildToken(buildClaim()))
		assert.Nil(t, claim)
		assert.NotNil(t, err)
	})
}