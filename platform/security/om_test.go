package security

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
)

const (
	userName = "test-User"
	password = "password"
)

var genericError = errors.New("error")

func buildUser(id string, password string) *domain.User {
	return &domain.User{
		ID:       id,
		UserName: userName,
		Email:    "",
		Password: []byte(password),
	}
}

func buildToken() string {
	return GenerateJWTToken(&domain.User{UserName: "test_user", ID: "id"})
}
