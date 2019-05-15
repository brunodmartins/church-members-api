package service

import (
	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/infra/mongo"
	"github.com/BrunoDM2943/church-members-api/member/repository"
)

type IMemberService interface {
	FindMembers(filters map[string]interface{}) ([]*entity.Membro, error)
	FindMembersByID(id entity.ID) (*entity.Membro, error)
	SaveMember(member *entity.Membro) (entity.ID, error)
}

type memberService struct {
	repo repository.IMemberRepository
}

func NewMemberService(r repository.IMemberRepository) *memberService {
	return &memberService{
		repo: r,
	}
}

func (s *memberService) FindMembers(filters map[string]interface{}) ([]*entity.Membro, error) {
	queryFilters := mongo.QueryFilters{}

	if sex := filters["sexo"]; sex != nil {
		queryFilters.AddFilter("pessoa.sexo", sex)
	}

	if active := filters["active"]; active != nil {
		queryFilters.AddFilter("active", active.(bool))
	}


	return s.repo.FindAll(queryFilters)
}

func (s *memberService) FindMembersByID(id entity.ID) (*entity.Membro, error) {
	return s.repo.FindByID(id)
}

func (s *memberService) SaveMember(member *entity.Membro) (entity.ID, error) {
	return s.repo.Insert(member)
}