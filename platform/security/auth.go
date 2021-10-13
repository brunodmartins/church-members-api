package security

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type Claim struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func newClaim(user *domain.User) *Claim {
	return &Claim{
		ID:       user.ID,
		UserName: user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: getExpirationTime(),
			Issuer:    "church-members-api",
		},
	}
}

func getExpirationTime() *jwt.NumericDate {
	hoursToExpire := viper.GetInt("security.token.expiration")
	return jwt.NewNumericDate(time.Now().Add(time.Duration(hoursToExpire) * time.Hour))
}
