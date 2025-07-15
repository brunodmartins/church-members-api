package church

import (
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	List() ([]*domain.Church, error)
	GetChurch(id string) (*domain.Church, error)
	GetChurchByAbbreviation(abbreviation string) (*domain.Church, error)
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

func (s churchService) GetChurchByAbbreviation(abbreviation string) (*domain.Church, error) {
	churches, err := s.repo.List()
	if err != nil {
		logrus.Errorf("error looking for church by abbreviation %s: %v", abbreviation, err)
		return nil, err
	}
	for _, ch := range churches {
		if ch.Abbreviation == abbreviation {
			return ch, nil
		}
	}
	err = apierrors.NewApiError(fmt.Sprintf("achurch for bbreviation %s not found", abbreviation), http.StatusNotFound)
	logrus.Infof("church for abbreviation %s not found", abbreviation)
	return nil, err
}
