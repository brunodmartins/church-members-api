package member

import (
	"errors"

	"github.com/BrunoDM2943/church-members-api/entity"
)

var (
	MemberNotFound, MemberError error
)

func init() {
	MemberNotFound = errors.New("Member not found")
	MemberError = errors.New("Member not found")
}

type Service interface {
	Reader
	Writer
}

type Reader interface {
	FindAll() ([]*entity.Membro, error)
	FindByID(id entity.ID) (*entity.Membro, error)
}

type Writer interface {
	Insert(membro *entity.Membro) (entity.ID, error)
}

//Repository repository interface
type Repository interface {
	Reader
	Writer
}

//UseCase use case interface
type UseCase interface {
	Reader
}
