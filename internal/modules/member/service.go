package member

import (
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	FindMembers(specification Specification) ([]*domain.Member, error)
	FindMembersByID(id string) (*domain.Member, error)
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

func (s *memberService) FindMembers(specification Specification) ([]*domain.Member, error) {
	return s.repo.FindAll(specification)
}

func (s *memberService) FindMembersByID(id string) (*domain.Member, error) {
	return s.repo.FindByID(id)
}

func (s *memberService) SaveMember(member *domain.Member) (string, error) {
	member.Active = true
	return s.repo.Insert(member)
}

func (s *memberService) ChangeStatus(ID string, status bool, reason string, date time.Time) error {
	err := s.repo.UpdateStatus(ID, status)
	if err == nil {
		return s.repo.GenerateStatusHistory(ID, status, reason, date)
	}
	return err
}
