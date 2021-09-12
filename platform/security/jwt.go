package security

import (
	"github.com/BrunoDM2943/church-members-api/platform/security/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GenerateJWTToken(claim *domain.Claim) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	result, _ := jwtToken.SignedString(getSecret())
	return result
}

func getSecret() []byte {
	return []byte(viper.GetString("security.token.secret"))
}

func getClaim(jwtToken string) (*domain.Claim, error) {
	token, err := parseJWT(jwtToken)
	if err != nil {
		logrus.Error("Error decrypting jwt token: ", err)
		return nil, err
	}
	result, _ := token.Claims.(*domain.Claim)
	return result, nil
}

func parseJWT(jwtToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		jwtToken,
		&domain.Claim{},
		func(token *jwt.Token) (interface{}, error) {
			return getSecret(), nil
		},
	)
}
