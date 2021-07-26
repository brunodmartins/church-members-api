package member

import (
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/repository"

	"github.com/BrunoDM2943/church-members-api/internal/constants/entity"
)

//go:generate mockgen -source=./memberService.go -destination=./mock/memberService_mock.go
type Service interface {
	FindMembers(filters map[string]interface{}) ([]*entity.Member, error)
	FindMembersByID(id string) (*entity.Member, error)
	SaveMember(member *entity.Member) (string, error)
	ChangeStatus(id string, status bool, reason string, date time.Time) error
}

type memberService struct {
	repo repository.MemberRepository
}

func NewMemberService(r repository.MemberRepository) *memberService {
	return &memberService{
		repo: r,
	}
}

func (s *memberService) FindMembers(filters map[string]interface{}) ([]*entity.Member, error) {
	queryFilters := repository.QueryFilters{}

	if sex := filters["gender"]; sex != nil {
		queryFilters.AddFilter("person.gender", sex)
	}

	if active := filters["active"]; active != nil {
		queryFilters.AddFilter("active", active.(bool))
	}

	if name := filters["name"]; name != nil {
		queryFilters.AddFilter("name", name)
	}

	return s.repo.FindAll(queryFilters)
}

func (s *memberService) FindMembersByID(id string) (*entity.Member, error) {
	return s.repo.FindByID(id)
}

func (s *memberService) SaveMember(member *entity.Member) (string, error) {
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
