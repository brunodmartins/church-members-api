package member

import (
	"github.com/BrunoDM2943/church-members-api/entity"
)

type MemberInMemoryRepository struct {
	dataSet []*entity.Membro
}

func NewMemberInMemoryRepository() *MemberInMemoryRepository {
	return &MemberInMemoryRepository{
		dataSet: make([]*entity.Membro, 0),
	}
}

func (repo *MemberInMemoryRepository) FindAll() ([]*entity.Membro, error) {
	return repo.dataSet, nil
}

func (repo *MemberInMemoryRepository) FindByID(id entity.ID) (*entity.Membro, error) {
	var result *entity.Membro
	for _, elem := range repo.dataSet {
		if elem.ID == id {
			result = elem
			break
		}
	}
	return result, nil
}
