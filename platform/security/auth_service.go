package security

import (
	"net/http"

	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/BrunoDM2943/church-members-api/platform/security/domain"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=./auth_service.go -destination=./mock/auth_service_mock.go
type Service interface {
	GenerateToken(username, password string) (string, error)
}

type authService struct {
	userRepository UserRepository
}

func NewAuthService(userRepository UserRepository) Service {
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
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", s.buildAuthError()
	}
	return GenerateJWTToken(domain.NewClaim(user)), nil
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
