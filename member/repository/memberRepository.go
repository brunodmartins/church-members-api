package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/infra/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//go:generate mockgen -source=./memberRepository.go -destination=./mock/memberRepository_mock.go
type IMemberRepository interface {
	FindAll(filters mongo.QueryFilters) ([]*entity.Member, error)
	FindByID(id entity.ID) (*entity.Member, error)
	Insert(member *entity.Member) (entity.ID, error)
	Search(text string) ([]*entity.Member, error)
	FindMonthBirthday(date time.Time) ([]*entity.Person, error)
	UpdateStatus(ID entity.ID, status bool) error
	GenerateStatusHistory(id entity.ID, status bool, reason string, date time.Time) error
}

type memberRepository struct {
	col        *mgo.Collection
	colHistory *mgo.Collection
}

var (
	MemberNotFound = errors.New("Member not found")
)

func NewMemberRepository(session *mgo.Session) *memberRepository {
	return &memberRepository{
		col:        session.DB("disciples").C("member"),
		colHistory: session.DB("disciples").C("member_history"),
	}
}

func (repo *memberRepository) FindAll(filters mongo.QueryFilters) ([]*entity.Member, error) {
	var result []*entity.Member
	err := repo.col.Find(filters).Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *memberRepository) FindByID(id entity.ID) (*entity.Member, error) {
	var result *entity.Member
	err := repo.col.FindId(bson.ObjectIdHex(id.String())).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			return nil, MemberNotFound
		}
		return nil, err
	}
	return result, nil
}

func (repo *memberRepository) Insert(member *entity.Member) (entity.ID, error) {
	member.ID = entity.NewID()
	return member.ID, repo.col.Insert(member)
}

func (repo *memberRepository) Search(text string) ([]*entity.Member, error) {
	var result []*entity.Member
	regex := bson.RegEx{fmt.Sprintf(".*%s*.", text), "i"}
	err := repo.col.Find(
		bson.M{
			"$or": []bson.M{
				{"person.firstName": regex},
				{"person.lastName": regex},
			},
		},
	).Select(bson.M{}).All(&result)
	return result, err
}

func (repo *memberRepository) FindMonthBirthday(date time.Time) ([]*entity.Person, error) {
	var result []*entity.Member
	var resultParsed []*entity.Person
	err := repo.col.Find(bson.M{
		"$expr": bson.M{
			"$eq": []interface{}{
				bson.M{
					"$month": "$person.dtNascimento",
				},
				date.Month(),
			},
		},
	}).All(&result)
	if err != nil {
		return nil, err
	}
	for _, member := range result {
		resultParsed = append(resultParsed, &member.Person)
	}
	return resultParsed, nil
}

func (repo *memberRepository) UpdateStatus(ID entity.ID, status bool) error {
	return repo.col.UpdateId(bson.ObjectIdHex(ID.String()), bson.M{
		"$set": bson.M{
			"active": status,
		}})
}

func (repo *memberRepository) GenerateStatusHistory(id entity.ID, status bool, reason string, date time.Time) error {
	return repo.colHistory.Insert(bson.M{
		"member_id":    id,
		"status":       status,
		"reason":       reason,
		"changed_date": date,
	})
}
