package user

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"net/http"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	SaveUser(ctx context.Context, user *domain.User) error
	SearchUser(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.User, error)
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s userService) SaveUser(ctx context.Context, user *domain.User) error {
	user.ChurchID = domain.GetChurchID(ctx)
	if err := s.checkUserExist(ctx, user.UserName); err != nil {
		return err
	}
	return s.repository.SaveUser(ctx, user)
}

func (s userService) SearchUser(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.User, error) {
	return s.repository.SearchUser(ctx, specification)
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
