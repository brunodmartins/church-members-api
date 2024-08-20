package user

import (
	"context"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/services/email"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/brunodmartins/church-members-api/platform/i18n"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"net/http"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	SaveUser(ctx context.Context, user *domain.User) error
	SearchUser(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.User, error)
	ConfirmEmail(ctx context.Context, userName string, token string) error
	SendConfirmEmail(ctx context.Context, user *domain.User) error
}

type userService struct {
	repository Repository
	email      email.Service
}

func NewService(repository Repository, email email.Service) Service {
	return &userService{repository: repository, email: email}
}

func (s userService) SaveUser(ctx context.Context, user *domain.User) error {
	if err := s.checkUserExist(ctx, user.UserName); err != nil {
		return err
	}
	user.ChurchID = domain.GetChurchID(ctx)
	if err := s.repository.SaveUser(ctx, user); err != nil {
		return err
	}
	if err := s.SendConfirmEmail(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s userService) SearchUser(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.User, error) {
	return s.repository.SearchUser(ctx, specification)
}

func (s userService) ConfirmEmail(ctx context.Context, userName string, token string) error {
	user, err := s.repository.FindUser(ctx, userName)
	if err != nil {
		return err
	}
	if user.ConfirmedEmail {
		return apierrors.NewApiError("email is already confirmed", http.StatusUnprocessableEntity)
	}
	if user.ConfirmationToken != token {
		return apierrors.NewApiError("The provided token does not match for the given user", http.StatusBadRequest)
	}
	user.ConfirmedEmail = true
	return s.repository.UpdateUser(ctx, user)
}

func (s userService) checkUserExist(ctx context.Context, userName string) error {
	user, err := s.repository.FindUser(ctx, userName)
	if err != nil {
		return err
	}
	if user != nil {
		return apierrors.NewApiError("User already exist", http.StatusBadRequest)
	}
	return nil
}

func (s userService) SendConfirmEmail(ctx context.Context, user *domain.User) error {
	data := email.NewConfirmEmailTemplateDTO(ctx)
	data.User = user.UserName
	data.Link = user.BuildConfirmationLink()
	return s.email.SendTemplateEmail(email.ConfirmEmailTemplate, data, i18n.GetMessage(ctx, "Emails.ConfirmEmail.Subject"), user.Email)
}
