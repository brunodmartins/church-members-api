package member

import (
	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/infra/mongo"
)

type MemberService struct {
	repo Repository
}

func NewMemberService(r Repository) *MemberService {
	return &MemberService{
		repo: r,
	}
}

func (s *MemberService) FindAll(filters map[string]interface{}) ([]*entity.Membro, error) {
	queryFilters := mongo.QueryFilters{}

	if sex := filters["sexo"]; sex != "" {
		queryFilters.AddFilter("pessoa.sexo", sex)
	}

	if active := filters["active"]; active != "" {
		queryFilters.AddFilter("active", active.(bool))
	}


	return s.repo.FindAll(queryFilters)
}

func (s *MemberService) FindByID(id entity.ID) (*entity.Membro, error) {
	return s.repo.FindByID(id)
}

func (s *MemberService) Insert(membro *entity.Membro) (entity.ID, error) {
	return s.repo.Insert(membro)
}