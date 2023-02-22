package church

import "github.com/brunodmartins/church-members-api/internal/constants/domain"

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	List() ([]*domain.Church, error)
	GetChurch(id string) (*domain.Church, error)
}

type churchService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &churchService{repo: repo}
}

func (s churchService) List() ([]*domain.Church, error) {
	return s.repo.List()
}

func (s churchService) GetChurch(id string) (*domain.Church, error) {
	return s.repo.GetByID(id)
}
