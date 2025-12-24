package church

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	List(ctx context.Context) ([]*domain.Church, error)
	GetChurch(ctx context.Context, id string) (*domain.Church, error)
	GetChurchByAbbreviation(ctx context.Context, abbreviation string) (*domain.Church, error)
	GetStatistics(ctx context.Context, id string) (*domain.ChurchStatistics, error)
}

type churchService struct {
	membersService member.Service
	repo           Repository
}

func NewService(memberService member.Service, repo Repository) Service {
	return &churchService{
		membersService: memberService,
		repo:           repo,
	}
}

func (s churchService) List(ctx context.Context) ([]*domain.Church, error) {
	return s.repo.List(ctx)
}

func (s churchService) GetChurch(ctx context.Context, id string) (*domain.Church, error) {
	if domain.GetChurchID(ctx) != id {
		return nil, apierrors.NewApiError("Not allowed to access other churches", http.StatusForbidden)
	}
	return s.repo.GetByID(ctx, id)
}

func (s churchService) GetChurchByAbbreviation(ctx context.Context, abbreviation string) (*domain.Church, error) {
	churches, err := s.repo.List(ctx)
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

func (s churchService) GetStatistics(ctx context.Context, id string) (*domain.ChurchStatistics, error) {
	if domain.GetChurchID(ctx) != id {
		return nil, apierrors.NewApiError("Not allowed to access other churches", http.StatusForbidden)
	}
	members, err := s.membersService.SearchMembers(ctx, member.OnlyActive())
	if err != nil {
		return nil, err
	}

	result := &domain.ChurchStatistics{
		TotalMembers:                 len(members),
		AgeDistribution:              make([]int, 0),
		TotalMembersByGender:         make(map[string]int),
		TotalMembersByClassification: make(map[string]int),
	}

	for _, m := range members {
		result.AgeDistribution = append(result.AgeDistribution, m.Person.Age())
		result.TotalMembersByGender[m.Person.Gender]++
		result.TotalMembersByClassification[m.Classification().String()]++
	}

	return result, nil
}
