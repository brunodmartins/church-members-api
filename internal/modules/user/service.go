package user

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"net/http"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	SaveUser(user *domain.User) error
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s userService) SaveUser(user *domain.User) error {
	if err := s.checkUserExist(user.UserName); err != nil {
		return err
	}
	return s.repository.SaveUser(user)
}

func (s userService) checkUserExist(userName string) error {
	user, err := s.repository.FindUser(userName)
	if err != nil {
		return err
	}
	if user != nil {
		return apierrors.NewApiError("User already exist", http.StatusBadRequest)
	}
	return nil
}
