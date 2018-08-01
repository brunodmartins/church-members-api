package member

import "github.com/BrunoDM2943/church-members-api/entity"

type Reader interface {
	FindAll() ([]*entity.Membro, error)
	FindByID(id entity.ID) (*entity.Membro, error)
}

//Repository repository interface
type Repository interface {
	Reader
}

//UseCase use case interface
type UseCase interface {
	Reader
}
