package utils

import (
	"time"

	"github.com/BrunoDM2943/church-members-api/entity"
)

//Repository repository interface
type Repository interface {
	FindMonthBirthday(date time.Time) ([]entity.Pessoa, error)
}

//UseCase use case interface
type UseCase interface {
	FindMonthBirthday(date time.Time) ([]entity.Pessoa, error)
}
