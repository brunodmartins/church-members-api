package member

import (
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"net/http"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	SearchMembers(specification Specification) ([]*domain.Member, error)
	GetMember(id string) (*domain.Member, error)
	SaveMember(member *domain.Member) (string, error)
	ChangeStatus(id string, status bool, reason string, date time.Time) error
}

type memberService struct {
	repo Repository
}

func NewMemberService(r Repository) *memberService {
	return &memberService{
		repo: r,
	}
}

func (s *memberService) SearchMembers(specification Specification) ([]*domain.Member, error) {
	return s.repo.FindAll(specification)
}

func (s *memberService) GetMember(id string) (*domain.Member, error) {
	if !domain.IsValidID(id) {
		return nil, apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	return s.repo.FindByID(id)
}

func (s *memberService) SaveMember(member *domain.Member) (string, error) {
	member.Active = true
	err := s.repo.Insert(member)
	return member.ID, err
}

func (s *memberService) ChangeStatus(ID string, status bool, reason string, date time.Time) error {
	member, err := s.GetMember(ID)
	if err != nil {
		return err
	}
	member.Active = status
	err = s.repo.UpdateStatus(member)
	if err == nil {
		return s.repo.GenerateStatusHistory(ID, status, reason, date)
	}
	return err
}
