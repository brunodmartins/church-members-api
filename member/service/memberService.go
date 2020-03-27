package service

import (
	"fmt"
	"time"

	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/infra/mongo"
	"github.com/BrunoDM2943/church-members-api/member/repository"
	"gopkg.in/mgo.v2/bson"
)

//go:generate mockgen -source=./memberService.go -destination=./mock/memberService_mock.go
type IMemberService interface {
	FindMembers(filters map[string]interface{}) ([]*entity.Membro, error)
	FindMembersByID(id entity.ID) (*entity.Membro, error)
	SaveMember(member *entity.Membro) (entity.ID, error)
	FindMonthBirthday(date time.Time) ([]*entity.Pessoa, error)
	ChangeStatus(id entity.ID, status bool, reason string) error
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

	if name := filters["name"]; name != nil {
		regex := bson.RegEx{fmt.Sprintf(".*%s*.", name), "i"}
		queryFilters.AddFilter("$or", []bson.M{
			{"pessoa.nome": regex},
			{"pessoa.sobrenome": regex},
		})
	}

	return s.repo.FindAll(queryFilters)
}

func (s *memberService) FindMembersByID(id entity.ID) (*entity.Membro, error) {
	return s.repo.FindByID(id)
}

func (s *memberService) SaveMember(member *entity.Membro) (entity.ID, error) {
	return s.repo.Insert(member)
}

func (s *memberService) FindMonthBirthday(month time.Time) ([]*entity.Pessoa, error) {
	return s.repo.FindMonthBirthday(month)
}

func (s *memberService) ChangeStatus(ID entity.ID, status bool, reason string) error {
	err := s.repo.UpdateStatus(ID, status)
	s.repo.GenerateStatusHistory(ID, status, reason, time.Now())
	return err
}
