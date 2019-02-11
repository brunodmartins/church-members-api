package utils

import (
	"time"

	"github.com/BrunoDM2943/church-members-api/entity"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UtilsRepository struct {
	col *mgo.Collection
}

func NewUtilsRepository(session *mgo.Session) *UtilsRepository {
	return &UtilsRepository{
		col: session.DB("disciples").C("Membro"),
	}
}

func (repo *UtilsRepository) FindMonthBirthday(date time.Time) ([]entity.Pessoa, error) {
	var result []*entity.Membro
	var resultParsed []entity.Pessoa
	/* db.Membro.find(
			{
				$expr: {
	    			$eq: [
						{
							$month: "$pessoa.dtNascimento"
						},
						1
					]
				}
			}
		)
	*/
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
		resultParsed = append(resultParsed, membro.Pessoa)
	}
	return resultParsed, nil
}
