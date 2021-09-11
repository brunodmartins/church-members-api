package security

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/platform/security/domain"
)

const (
	userName = "test-User"
	password = "password"
)

var genericError = errors.New("error")

func buildUser(id string) *domain.User {
	return &domain.User{
		ID:       id,
		UserName: userName,
		Email:    "",
		Password: password,
	}
}

func buildToken(claim *domain.Claim) string {
	return GenerateJWTToken(claim)
}

func buildClaim() *domain.Claim {
	return domain.NewClaim(&domain.User{UserName: "test_user", ID: "id"})
}