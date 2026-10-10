package user

import (
	"context"
	"net/http"
	"slices"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	SaveUser(ctx context.Context, user *domain.User) error
	FindUser(ctx context.Context, userName string) (*domain.User, error)
	SearchUser(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userName string) error
}

type userService struct {
	repository Repository
}

func (s userService) UpdateUser(ctx context.Context, user *domain.User) error {
	if !domain.GetUser(ctx).IsAdmin() {
		current, _ := s.repository.FindUser(ctx, user.UserName)
		if !slices.Equal(current.Roles, user.Roles) {
			return apierrors.NewApiError("User does not have required role", http.StatusForbidden)
		}
	}
	return s.repository.UpdateUser(ctx, user)
}

func (s userService) DeleteUser(ctx context.Context, userName string) error {
	if !domain.GetUser(ctx).IsAdmin() {
		return apierrors.NewApiError("User does not have required role", http.StatusForbidden)
	}
	if domain.GetUser(ctx).UserName == userName {
		return apierrors.NewApiError("User cannot delete themselves", http.StatusForbidden)
	}
	user, err := s.repository.FindUser(ctx, userName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return apierrors.NewApiError("User not found", http.StatusNotFound)
		}
		logrus.Errorf("Error finding user: %v", err)
		return err
	}
	err = s.repository.DeleteUser(ctx, user)
	if err != nil {
		logrus.Errorf("Error deleting user: %v", err)
		return err
	}
	logrus.Infof("User %s deleted successfully", userName)
	return nil
}

func (s userService) FindUser(ctx context.Context, userName string) (*domain.User, error) {
	return s.repository.FindUser(ctx, userName)
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s userService) SaveUser(ctx context.Context, user *domain.User) error {
	if err := s.checkUserExist(ctx, user.UserName); err != nil {
		return err
	}
	user.ChurchID = domain.GetChurchID(ctx)
	return s.repository.SaveUser(ctx, user)
}

func (s userService) SearchUser(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.User, error) {
	return s.repository.SearchUser(ctx, specification)
}

func (s userService) checkUserExist(ctx context.Context, userName string) error {
	user, err := s.repository.FindUser(ctx, userName)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	if user != nil {
		return apierrors.NewApiError("User already exist", http.StatusBadRequest)
	}
	return nil
}
