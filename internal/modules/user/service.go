package user

import "github.com/BrunoDM2943/church-members-api/internal/constants/domain"

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
	return s.repository.SaveUser(user)
}
