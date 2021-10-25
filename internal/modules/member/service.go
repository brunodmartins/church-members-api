package member

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"net/http"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	SearchMembers(ctx context.Context, querySpecification wrapper.QuerySpecification, postSpecification ...Specification) ([]*domain.Member, error)
	GetMember(ctx context.Context, id string) (*domain.Member, error)
	SaveMember(ctx context.Context, member *domain.Member) (string, error)
	ChangeStatus(ctx context.Context, id string, status bool, reason string, date time.Time) error
}

type memberService struct {
	repo Repository
}

func NewMemberService(r Repository) Service {
	return &memberService{
		repo: r,
	}
}

func (s *memberService) SearchMembers(ctx context.Context, querySpecification wrapper.QuerySpecification, postSpecification ...Specification) ([]*domain.Member, error) {
	members, err := s.repo.FindAll(ctx, querySpecification)
	if err != nil {
		return nil, err
	}
	if len(postSpecification) != 0 {
		return applySpecifications(members, postSpecification), nil
	}
	return members, nil
}

func (s *memberService) GetMember(ctx context.Context, id string) (*domain.Member, error) {
	if !domain.IsValidID(id) {
		return nil, apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	return s.repo.FindByID(ctx, id)
}

func (s *memberService) SaveMember(ctx context.Context, member *domain.Member) (string, error) {
	member.Active = true
	err := s.repo.Insert(ctx, member)
	return member.ID, err
}

func (s *memberService) ChangeStatus(ctx context.Context, ID string, status bool, reason string, date time.Time) error {
	member, err := s.GetMember(ctx, ID)
	if err != nil {
		return err
	}
	member.Active = status
	err = s.repo.UpdateStatus(ctx, member)
	if err == nil {
		return s.repo.GenerateStatusHistory(ID, status, reason, date)
	}
	return err
}
