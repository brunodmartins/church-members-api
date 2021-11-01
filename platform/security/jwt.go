package security

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GenerateJWTToken(user *domain.User) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaim(user))
	result, _ := jwtToken.SignedString(getSecret())
	return result
}

func getSecret() []byte {
	return []byte(viper.GetString("security.token.secret"))
}

func getClaim(jwtToken string) (*Claim, error) {
	token, err := parseJWT(jwtToken)
	if err != nil {
		logrus.Error("Error decrypting jwt token: ", err)
		return nil, err
	}
	result, _ := token.Claims.(*Claim)
	return result, nil
}

func parseJWT(jwtToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		jwtToken,
		&Claim{},
		func(token *jwt.Token) (interface{}, error) {
			return getSecret(), nil
		},
	)
}
