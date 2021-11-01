package security

import (
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type Claim struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Church   *domain.Church
	jwt.RegisteredClaims
}

func newClaim(user *domain.User) *Claim {
	return &Claim{
		ID:       user.ID,
		UserName: user.UserName,
		Church:   user.Church,
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
