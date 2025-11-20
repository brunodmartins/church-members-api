package security

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/services/email"
	"github.com/brunodmartins/church-members-api/platform/i18n"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/church"
	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/brunodmartins/church-members-api/platform/crypto"

	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
)

//go:generate mockgen -source=./auth_service.go -destination=./mock/auth_service_mock.go
type Service interface {
	IdentifyChurch(ctx context.Context, churchAbbreviation, churchID string) (*domain.Church, error)
	GenerateToken(church *domain.Church, username, password string) (string, error)
	SendConfirmEmail(ctx context.Context, user *domain.User) error
	ConfirmEmail(ctx context.Context, userName string) error
}

type authService struct {
	userService   user.Service
	churchService church.Service
	email         email.Service
}

func NewAuthService(userService user.Service, churchService church.Service, emailService email.Service) Service {
	return &authService{
		userService:   userService,
		churchService: churchService,
		email:         emailService,
	}
}

func (s *authService) GenerateToken(church *domain.Church, username, password string) (string, error) {
	ctx := context.WithValue(context.Background(), "church", church)
	user, err := s.userService.FindUser(ctx, username)
	if err != nil {
		logrus.Errorf("Error getting user %s: %v", username, err)
		return "", err
	}
	err = crypto.IsSamePassword(user.Password, password)
	if err != nil {
		logrus.Infof("Provided password user %s not equal", username)
		return "", apierrors.NewApiError("Invalid password", http.StatusUnauthorized)
	}
	user.Church = church
	if !user.ConfirmedEmail {
		logrus.Infof("User %s has not confirmed the email, confirmation mail sent", username)
		if err = s.SendConfirmEmail(ctx, user); err != nil {
			return "", err
		}
		return "", apierrors.NewApiError("User requires email confirmation", http.StatusUnprocessableEntity)
	}
	return GenerateJWTToken(user), nil
}

func (s *authService) SendConfirmEmail(ctx context.Context, user *domain.User) error {
	data := email.NewConfirmEmailTemplateDTO(ctx)
	data.User = user.UserName
	data.Link = buildConfirmationLink(GenerateJWTToken(user))

	return s.email.SendTemplateEmail(email.ConfirmEmailTemplate, data, i18n.GetMessage(ctx, "Emails.ConfirmEmail.Subject"), user.Email)
}

func buildConfirmationLink(accessToken string) string {
	return fmt.Sprintf("%s/users/confirm?accessToken=%s", viper.GetString("email.confirm.url"), accessToken)
}

func (s *authService) ConfirmEmail(ctx context.Context, userName string) error {
	user, err := s.userService.FindUser(ctx, userName)
	if err != nil {
		return err
	}
	if user.ConfirmedEmail {
		return apierrors.NewApiError("email is already confirmed", http.StatusUnprocessableEntity)
	}
	user.ConfirmedEmail = true
	return s.userService.UpdateUser(ctx, user)
}

func (s *authService) IdentifyChurch(ctx context.Context, churchAbbreviation, churchID string) (*domain.Church, error) {
	if churchAbbreviation == "" && churchID == "" {
		return nil, apierrors.NewApiError("Church ID or abbreviation must be provided", http.StatusBadRequest)
	}

	if churchAbbreviation != "" {
		return s.churchService.GetChurchByAbbreviation(ctx, churchAbbreviation)
	} else {
		if !domain.IsValidID(churchID) {
			return nil, apierrors.NewApiError(fmt.Sprintf("Invalid Church ID %s", churchID), http.StatusBadRequest)
		}
		return s.churchService.GetChurch(ctx, churchID)
	}
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
