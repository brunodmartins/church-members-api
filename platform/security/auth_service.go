package security

import (
	"context"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/church"
	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/brunodmartins/church-members-api/platform/crypto"

	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
)

//go:generate mockgen -source=./auth_service.go -destination=./mock/auth_service_mock.go
type Service interface {
	GenerateToken(churchID, username, password string) (string, error)
}

type authService struct {
	userRepository   user.Repository
	churchRepository church.Repository
}

func NewAuthService(userRepository user.Repository, churchRepository church.Repository) Service {
	return &authService{
		userRepository:   userRepository,
		churchRepository: churchRepository,
	}
}

func (s *authService) GenerateToken(churchID, username, password string) (string, error) {
	church, err := s.churchRepository.GetByID(churchID)
	if err != nil {
		return "", s.buildAuthError()
	}
	user, err := s.userRepository.FindUser(context.WithValue(context.Background(), "user", &domain.User{
		Church: church,
	}), username)
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
	user.Church = church
	return GenerateJWTToken(user), nil
}

func (s *authService) buildAuthError() apierrors.Error {
	return apierrors.NewApiError("User not found. Check information.", http.StatusNotFound)
}

func GetClaim(token string) (bool, *Claim) {
	claim, err := getClaim(token)
	if err != nil {
		return false, claim
	}
	return true, claim
}

func AddClaimToContext(claim *Claim, ctx context.Context) context.Context {
	return context.WithValue(ctx, "user", &domain.User{
		ID:       claim.ID,
		UserName: claim.UserName,
		Church:   claim.Church,
	})
}
