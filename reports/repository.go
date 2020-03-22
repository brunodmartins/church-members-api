package reports

import (
	"time"

	"github.com/BrunoDM2943/church-members-api/entity"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository_mock.go
type ReportRepository interface {
	FindMonthBirthday(date time.Time) ([]*entity.Pessoa, error)
}

type reportRepositoryMongo struct {
	col *mgo.Collection
}

func NewReportRepositoryMongo(session *mgo.Session) ReportRepository {
	return reportRepositoryMongo{
		col: session.DB("disciples").C("Membro"),
	}
}

func (repo reportRepositoryMongo) FindMonthBirthday(date time.Time) ([]*entity.Pessoa, error) {
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
