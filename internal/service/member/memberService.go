package member

import (
	"fmt"
	"github.com/BrunoDM2943/church-members-api/internal/repository"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/BrunoDM2943/church-members-api/internal/storage/mongo"
	"gopkg.in/mgo.v2/bson"
)

//go:generate mockgen -source=./memberService.go -destination=./mock/memberService_mock.go
type Service interface {
	FindMembers(filters map[string]interface{}) ([]*model.Member, error)
	FindMembersByID(id model.ID) (*model.Member, error)
	SaveMember(member *model.Member) (model.ID, error)
	FindMonthBirthday(date time.Time) ([]*model.Person, error)
	ChangeStatus(id model.ID, status bool, reason string) error
}

type memberService struct {
	repo repository.MemberRepository
}

func NewMemberService(r repository.MemberRepository) *memberService {
	return &memberService{
		repo: r,
	}
}

func (s *memberService) FindMembers(filters map[string]interface{}) ([]*model.Member, error) {
	queryFilters := mongo.QueryFilters{}

	if sex := filters["gender"]; sex != nil {
		queryFilters.AddFilter("person.gender", sex)
	}

	if active := filters["active"]; active != nil {
		queryFilters.AddFilter("active", active.(bool))
	}

	if name := filters["name"]; name != nil {
		regex := bson.RegEx{fmt.Sprintf(".*%s*.", name), "i"}
		queryFilters.AddFilter("$or", []bson.M{
			{"person.firstName": regex},
			{"person.lastName": regex},
		})
	}

	return s.repo.FindAll(queryFilters)
}

func (s *memberService) FindMembersByID(id model.ID) (*model.Member, error) {
	return s.repo.FindByID(id)
}

func (s *memberService) SaveMember(member *model.Member) (model.ID, error) {
	return s.repo.Insert(member)
}

func (s *memberService) FindMonthBirthday(month time.Time) ([]*model.Person, error) {
	return s.repo.FindMonthBirthday(month)
}

func (s *memberService) ChangeStatus(ID model.ID, status bool, reason string) error {
	err := s.repo.UpdateStatus(ID, status)
	s.repo.GenerateStatusHistory(ID, status, reason, time.Now())
	return err
}
