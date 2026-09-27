package security

import (
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/enum"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type Claim struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Church   *domain.Church
	Roles    []string `json:"roles"`
	Role     string   `json:"role"`
	jwt.RegisteredClaims
}

func newClaim(user *domain.User) *Claim {
	return &Claim{
		ID:       user.ID,
		UserName: user.UserName,
		Church:   user.Church,
		Roles:    user.Roles,
		Role:     roleToString(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: getExpirationTime(),
			Issuer:    "church-members-api",
		},
	}
}

func roleToString(userRole enum.Role) string {
	if userRole < 0 || int(userRole) > 1 {
		return ""
	}
	return userRole.String()
}

func getExpirationTime() *jwt.NumericDate {
	hoursToExpire := viper.GetInt("security.token.expiration")
	return jwt.NewNumericDate(time.Now().Add(time.Duration(hoursToExpire) * time.Hour))
}
