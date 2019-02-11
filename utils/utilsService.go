package utils

import (
	"time"

	"github.com/BrunoDM2943/church-members-api/entity"
)

type UtilsService struct {
	repo Repository
}

func NewUtilsService(repo Repository) *UtilsService {
	return &UtilsService{repo}
}

func (service *UtilsService) FindMonthBirthday(month time.Time) ([]entity.Pessoa, error) {
	return service.repo.FindMonthBirthday(month)
}
