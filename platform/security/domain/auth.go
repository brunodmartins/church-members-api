package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
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
	jwt.StandardClaims
}

func NewClaim(user *User) *Claim {
	return &Claim{
		ID: user.ID,
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: GetExpirationTime(),
			Issuer:    "church-members-api",
		},
	}
}

func GetExpirationTime() int64 {
	hoursToExpire := viper.GetInt("security.token.expiration")
	return time.Now().UTC().Add(time.Duration(hoursToExpire) * time.Hour).Unix()
}