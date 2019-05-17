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

type IMemberRepository interface {
	FindAll(filters mongo.QueryFilters) ([]*entity.Membro, error)
	FindByID(id entity.ID) (*entity.Membro, error)
	Insert(membro *entity.Membro) (entity.ID, error)
	Search(text string) ([]*entity.Membro, error)
	FindMonthBirthday(date time.Time) ([]*entity.Pessoa, error)
}

type memberRepository struct {
	col *mgo.Collection
}

var (
	MemberNotFound = errors.New("Member not found")
)

func NewMemberRepository(session *mgo.Session) *memberRepository {
	return &memberRepository{
		col: session.DB("disciples").C("Membro"),
	}
}

func (repo *memberRepository) FindAll(filters mongo.QueryFilters) ([]*entity.Membro, error) {
	var result []*entity.Membro
	err := repo.col.Find(filters).Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *memberRepository) FindByID(id entity.ID) (*entity.Membro, error) {
	var result *entity.Membro
	err := repo.col.FindId(bson.ObjectIdHex(id.String())).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			return nil, MemberNotFound
		}
		return nil, err
	}
	return result, nil
}

func (repo *memberRepository) Insert(membro *entity.Membro) (entity.ID, error) {
	membro.ID = entity.NewID()
	return membro.ID, repo.col.Insert(membro)
}

func (repo *memberRepository) Search(text string) ([]*entity.Membro, error) {
	var result []*entity.Membro
	regex := bson.RegEx{fmt.Sprintf(".*%s*.", text), "i"}
	err := repo.col.Find(
		bson.M{
			"$or": []bson.M{
				{"pessoa.nome": regex},
				{"pessoa.sobrenome": regex},
			},
		},
	).Select(bson.M{}).All(&result)
	return result, err
}

func (repo *memberRepository) FindMonthBirthday(date time.Time) ([]*entity.Pessoa, error) {
	var result []*entity.Membro
	var resultParsed []*entity.Pessoa
	err := repo.col.Find(bson.M{
		"$expr": bson.M{
			"$eq": []interface{}{
				bson.M{
					"$month": "$pessoa.dtNascimento",
				},
				date.Month(),
			},
		},
	}).All(&result)
	if err != nil {
		return nil, err
	}
	for _, membro := range result {
		resultParsed = append(resultParsed, &membro.Pessoa)
	}
	return resultParsed, nil
}
