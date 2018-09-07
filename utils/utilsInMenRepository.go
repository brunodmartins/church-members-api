package utils

import (
	"errors"
	"time"

	"github.com/BrunoDM2943/church-members-api/entity"
)

type UtilsInMenRepository struct {
}

func NewUtilsInMenRepository() *UtilsInMenRepository {
	return &UtilsInMenRepository{}
}

func (repo *UtilsInMenRepository) FindMonthBirthday(date time.Time) ([]entity.Pessoa, error) {
	if date.Month() == time.January {
		return nil, errors.New("Mes invalido")
	} else {
		p := entity.Pessoa{
			Nome: "Bruno",
		}
		return []entity.Pessoa{p}, nil
	}

}
