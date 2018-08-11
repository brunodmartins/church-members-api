package member

import (
	"github.com/BrunoDM2943/church-members-api/entity"
)

type MemberInMemoryRepository struct {
	dataset []*entity.Membro
}

func NewMemberInMemoryRepository() *MemberInMemoryRepository {
	return &MemberInMemoryRepository{
		dataset: make([]*entity.Membro, 0),
	}
}

func (repo *MemberInMemoryRepository) FindAll() ([]*entity.Membro, error) {
	return repo.dataset, nil
}

func (repo *MemberInMemoryRepository) InsertMember(membro *entity.Membro) (entity.ID, error) {
	membro.ID = entity.NewID()
	repo.dataset = append(repo.dataset, membro)
	return membro.ID, nil
}

func (repo *MemberInMemoryRepository) FindByID(id entity.ID) (*entity.Membro, error) {
	if id.String() == "" {
		return nil, MemberError
	}
	var result *entity.Membro
	for _, elem := range repo.dataset {
		if elem.ID == id {
			result = elem
			break
		}
	}
	if result == nil {
		return nil, MemberNotFound
	}
	return result, nil
}
