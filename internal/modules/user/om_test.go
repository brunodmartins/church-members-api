package user

import (
	"errors"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/role"
)

const (
	userName          = "test-User"
	password          = "password"
	confirmationToken = "token"
)

var genericError = errors.New("error")

func buildUser(id string, password string) *domain.User {
	return &domain.User{
		ID:                id,
		UserName:          userName,
		Email:             "",
		Password:          []byte(password),
		Role:              role.USER,
		ConfirmedEmail:    false,
		ConfirmationToken: confirmationToken,
	}
}
