package confirmation

import (
	"errors"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

const confirmationTokenPurpose = "email_confirmation"

type confirmationClaim struct {
	ChurchID          string `json:"church_id"`
	ConfirmationToken string `json:"confirmation_token"`
	Purpose           string `json:"purpose"`
	jwt.RegisteredClaims
}

func GenerateToken(user *domain.User) (string, error) {
	if user == nil {
		return "", errors.New("user is required")
	}
	secret := getSecret()
	if len(secret) == 0 {
		return "", errors.New("token signing secret is required")
	}

	claims := confirmationClaim{
		ChurchID:          user.ChurchID,
		ConfirmationToken: user.ConfirmationToken,
		Purpose:           confirmationTokenPurpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    "church-members-api",
			Subject:   user.ID,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func ValidateToken(token string, user *domain.User) error {
	if user == nil {
		return errors.New("user is required")
	}
	secret := getSecret()
	if len(secret) == 0 {
		return errors.New("token signing secret is required")
	}

	claims := new(confirmationClaim)
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return err
	}
	if !parsed.Valid ||
		claims.Issuer != "church-members-api" ||
		claims.Subject != user.ID ||
		claims.ChurchID != user.ChurchID ||
		claims.ConfirmationToken != user.ConfirmationToken ||
		claims.Purpose != confirmationTokenPurpose {
		return errors.New("invalid confirmation token")
	}
	return nil
}

func getSecret() []byte {
	return []byte(viper.GetString("security.token.secret"))
}
