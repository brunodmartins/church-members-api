package church

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/services/storage"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/sirupsen/logrus"
)

const churchLogoFileExtension = ".png"

func churchLogoFileName(abbreviation string) string {
	return abbreviation + churchLogoFileExtension
}

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	List(ctx context.Context) ([]*domain.Church, error)
	GetChurch(ctx context.Context, id string) (*domain.Church, error)
	GetChurchByAbbreviation(ctx context.Context, abbreviation string) (*domain.Church, error)
	GetStatistics(ctx context.Context, id string) (*domain.ChurchStatistics, error)
	UpdateChurch(ctx context.Context, church *domain.Church) error
	UpdateChurchLogo(ctx context.Context, id string, data []byte) error
}

type churchService struct {
	membersService member.Service
	repo           Repository
	storageService storage.Service
}

func NewService(memberService member.Service, repo Repository, storageService storage.Service) Service {
	return &churchService{
		membersService: memberService,
		repo:           repo,
		storageService: storageService,
	}
}

func (s churchService) List(ctx context.Context) ([]*domain.Church, error) {
	return s.repo.List(ctx)
}

func (s churchService) GetChurch(ctx context.Context, id string) (*domain.Church, error) {
	if domain.GetChurchID(ctx) != id {
		return nil, apierrors.NewApiError("Not allowed to access other churches", http.StatusForbidden)
	}
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	logoFileName := churchLogoFileName(result.Abbreviation)
	result.Logo, err = s.storageService.GetFileURL(ctx, logoFileName)
	if err != nil {
		return nil, err
	}
	return result, nil
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

func (s churchService) UpdateChurch(ctx context.Context, updatedChurch *domain.Church) error {
	currentChurch, err := s.getChurchForUpdate(ctx, updatedChurch.ID)
	if err != nil {
		return err
	}

	currentChurch.Name = updatedChurch.Name
	currentChurch.Language = updatedChurch.Language
	currentChurch.Email = updatedChurch.Email

	return s.repo.Update(ctx, currentChurch)
}

func (s churchService) UpdateChurchLogo(ctx context.Context, id string, data []byte) error {
	currentChurch, err := s.getChurchForUpdate(ctx, id)
	if err != nil {
		return err
	}

	logoFileName := churchLogoFileName(currentChurch.Abbreviation)
	if err := s.storageService.SaveFile(ctx, logoFileName, data); err != nil {
		return err
	}
	return nil
}

func (s churchService) getChurchForUpdate(ctx context.Context, id string) (*domain.Church, error) {
	currentUser := domain.GetUser(ctx)
	if currentUser == nil || !currentUser.IsAdmin() {
		return nil, apierrors.NewApiError("User does not have required role", http.StatusForbidden)
	}
	if currentUser.Church == nil || currentUser.Church.ID == "" {
		return nil, apierrors.NewApiError("Not allowed to access other churches", http.StatusForbidden)
	}
	if currentUser.Church.ID != id {
		return nil, apierrors.NewApiError("Not allowed to access other churches", http.StatusForbidden)
	}

	currentChurch, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return currentChurch, nil
}
