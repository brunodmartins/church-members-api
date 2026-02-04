package participant

import (
	"context"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	CreateParticipant(ctx context.Context, participant *domain.Participant) (string, error)
	GetParticipant(ctx context.Context, id string) (*domain.Participant, error)
	UpdateParticipant(ctx context.Context, participant *domain.Participant) error
	DeleteParticipant(ctx context.Context, id string) error
	SearchParticipant(ctx context.Context, querySpecification wrapper.QuerySpecification, postSpecification ...Specification) ([]*domain.Participant, error)
}

type participantService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &participantService{repo: repo}
}

func (s participantService) CreateParticipant(ctx context.Context, participant *domain.Participant) (string, error) {
	participant.ChurchID = domain.GetChurchID(ctx)
	err := s.repo.Insert(ctx, participant)
	return participant.ID, err
}

func (s participantService) GetParticipant(ctx context.Context, id string) (*domain.Participant, error) {
	if !domain.IsValidID(id) {
		return nil, apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	return s.repo.FindByID(ctx, id)
}

func (s participantService) UpdateParticipant(ctx context.Context, updateParticipant *domain.Participant) error {
	p, err := s.GetParticipant(ctx, updateParticipant.ID)
	if err != nil {
		return err
	}
	p.Name = updateParticipant.Name
	p.Filiation = updateParticipant.Filiation
	p.BirthDate = updateParticipant.BirthDate
	p.Observation = updateParticipant.Observation
	p.CellPhone = updateParticipant.CellPhone
	return s.repo.Update(ctx, p)
}

func (s participantService) DeleteParticipant(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s participantService) SearchParticipant(ctx context.Context, querySpecification wrapper.QuerySpecification, postSpecification ...Specification) ([]*domain.Participant, error) {
	participants, err := s.repo.FindAll(ctx, querySpecification)
	if err != nil {
		return nil, err
	}
	if len(postSpecification) != 0 {
		return applySpecifications(participants, postSpecification), nil
	}
	return participants, nil
}
