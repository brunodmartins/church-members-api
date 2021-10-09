package security

import (
	"github.com/BrunoDM2943/church-members-api/internal/modules/user"
	"github.com/BrunoDM2943/church-members-api/platform/crypto"
	"net/http"

	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
)

//go:generate mockgen -source=./auth_service.go -destination=./mock/auth_service_mock.go
type Service interface {
	GenerateToken(username, password string) (string, error)
}

type authService struct {
	userRepository user.Repository
}

func NewAuthService(userRepository user.Repository) Service {
	return &authService{userRepository: userRepository}
}

func (s *authService) GenerateToken(username, password string) (string, error) {
	user, err := s.userRepository.FindUser(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", s.buildAuthError()
	}
	err = crypto.IsSamePassword(user.Password, password)
	if err != nil {
		return "", s.buildAuthError()
	}
	return GenerateJWTToken(user), nil
}

func (s *authService) buildAuthError() apierrors.Error {
	return apierrors.NewApiError("User not found. Check information.", http.StatusNotFound)
}

func IsValidToken(token string) bool {
	_, err := getClaim(token)
	if err != nil {
		return false
	}
	return true
}
