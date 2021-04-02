package repository

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/BrunoDM2943/church-members-api/internal/storage/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//go:generate mockgen -source=./reportsRepository.go -destination=./mock/reportsRepository_mock.go
type ReportRepository interface {
	FindMembersActive() ([]*model.Member, error)
	FindMembersActiveAndMarried() ([]*model.Member, error)
}

type reportRepositoryImpl struct {
	col *mgo.Collection
}

func NewReportRepository(session *mgo.Session) ReportRepository {
	return &reportRepositoryImpl{
		col: session.DB("disciples").C("member"),
	}
}

func (repo reportRepositoryImpl) FindMembersActive() ([]*model.Member, error) {
	var result []*model.Member
	filters := mongo.QueryFilters{}
	filters.AddFilter("active", true)
	err := repo.col.Find(filters).Sort("person.firstName", "person.lastName").Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (repo reportRepositoryImpl) FindMembersActiveAndMarried() ([]*model.Member, error) {
	var result []*model.Member
	filters := mongo.QueryFilters{}
	filters.AddFilter("active", true)
	filters.AddFilter("person.marriageDate", bson.M{"$exists": true})
	err := repo.col.Find(filters).Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
