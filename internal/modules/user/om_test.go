package user

import (
	"errors"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

const (
	userName          = "test-User"
	password          = "password"
	confirmationToken = "token"
)

var genericError = errors.New("error")

func buildUser(id string, password string, roles ...string) *domain.User {
	return &domain.User{
		ID:             id,
		UserName:       userName,
		Email:          "",
		Password:       []byte(password),
		ConfirmedEmail: false,
		Roles:          roles,
	}
}

func buildAdminUser() *domain.User {
	return &domain.User{
		ID:             "admin-id",
		UserName:       "admin-User",
		Email:          "",
		Password:       []byte("password"),
		ConfirmedEmail: false,
		Roles:          []string{"admin"},
	}
}
