package security

import (
	"context"
	"github.com/sirupsen/logrus"
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
	userService   user.Service
	churchService church.Service
}

func NewAuthService(userService user.Service, churchService church.Service) Service {
	return &authService{
		userService:   userService,
		churchService: churchService,
	}
}

func (s *authService) GenerateToken(churchID, username, password string) (string, error) {
	church, err := s.churchService.GetChurch(churchID)
	if err != nil {
		logrus.Errorf("Error getting church %s: %v", churchID, err)
		return "", s.buildAuthError()
	}
	ctx := context.WithValue(context.Background(), "church", church)
	users, err := s.userService.SearchUser(ctx, user.WithUserName(username))
	if err != nil {
		logrus.Errorf("Error getting user %s: %v", username, err)
		return "", err
	}

	if len(users) == 0 {
		logrus.Infof("User %s not found", username)
		return "", s.buildAuthError()
	}
	user := users[0]
	err = crypto.IsSamePassword(user.Password, password)
	if err != nil {
		logrus.Infof("Provided password user %s not equal", username)
		return "", s.buildAuthError()
	}
	if !user.ConfirmedEmail {
		logrus.Infof("User %s has not confirmed the email, confirmation mail sent", username)
		_ = s.userService.SendConfirmEmail(ctx, user)
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
