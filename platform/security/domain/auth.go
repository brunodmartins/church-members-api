package domain

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claim struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func NewClaim(user *User) *Claim {
	return &Claim{
		ID:       user.ID,
		UserName: user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: GetExpirationTime(),
			Issuer:    "church-members-api",
		},
	}
}

func GetExpirationTime() *jwt.NumericDate {
	hoursToExpire := viper.GetInt("security.token.expiration")
	return jwt.NewNumericDate(time.Now().Add(time.Duration(hoursToExpire) * time.Hour))
}
