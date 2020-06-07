package reports

import (
	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/infra/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//go:generate mockgen -source=./reportsRepository.go -destination=./mock/reportsRepository_mock.go
type ReportRepository interface {
	FindMembersActive() ([]*entity.Member, error)
	FindMembersActiveAndMarried() ([]*entity.Member, error)
}

type reportRepositoryImpl struct {
	col *mgo.Collection
}

func NewReportRepository(session *mgo.Session) ReportRepository {
	return &reportRepositoryImpl{
		col: session.DB("disciples").C("Member"),
	}
}

func (repo reportRepositoryImpl) FindMembersActive() ([]*entity.Member, error) {
	var result []*entity.Member
	filters := mongo.QueryFilters{}
	filters.AddFilter("active", true)
	err := repo.col.Find(filters).Sort("pessoa.nome", "pessoa.sobrenome").Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (repo reportRepositoryImpl) FindMembersActiveAndMarried() ([]*entity.Member, error) {
	var result []*entity.Member
	filters := mongo.QueryFilters{}
	filters.AddFilter("active", true)
	filters.AddFilter("pessoa.dtCasamento", bson.M{"$exists": true})
	err := repo.col.Find(filters).Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
