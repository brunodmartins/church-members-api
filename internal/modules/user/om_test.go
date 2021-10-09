package user

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum/role"
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
		Role:     role.READ_ONLY,
	}
}

