package reports

import (
	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/infra/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//go:generate mockgen -source=./reportsRepository.go -destination=./mock/reportsRepository_mock.go
type ReportRepository interface {
	FindMembersActive() ([]*entity.Membro, error)
	FindMembersActiveAndMarried() ([]*entity.Membro, error)
}

type reportRepositoryImpl struct {
	col *mgo.Collection
}

func NewReportRepository(session *mgo.Session) ReportRepository {
	return &reportRepositoryImpl{
		col: session.DB("disciples").C("Membro"),
	}
}

func (repo reportRepositoryImpl) FindMembersActive() ([]*entity.Membro, error) {
	var result []*entity.Membro
	filters := mongo.QueryFilters{}
	filters.AddFilter("active", true)
	err := repo.col.Find(filters).Sort("pessoa.nome", "pessoa.sobrenome").Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (repo reportRepositoryImpl) FindMembersActiveAndMarried() ([]*entity.Membro, error) {
	var result []*entity.Membro
	filters := mongo.QueryFilters{}
	filters.AddFilter("active", true)
	filters.AddFilter("pessoa.dtCasamento", bson.M{"$exists": true})
	err := repo.col.Find(filters).Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
