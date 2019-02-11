package member

import (
	"github.com/BrunoDM2943/church-members-api/entity"
	"strings"
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

func (repo *MemberInMemoryRepository) Insert(membro *entity.Membro) (entity.ID, error) {
	membro.ID = entity.NewID()
	repo.dataSet = append(repo.dataSet, membro)
	return membro.ID, nil
}

func (repo *MemberInMemoryRepository) FindByID(id entity.ID) (*entity.Membro, error) {
	if id.String() == "" {
		return nil, MemberError
	}
	var result *entity.Membro
	for _, elem := range repo.dataSet {
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

func (repo *MemberInMemoryRepository) Search(text string) ([]*entity.Membro, error) {
	result := make([]*entity.Membro, 0)
	for _, m := range repo.dataSet {
		containsName := strings.Contains(m.Pessoa.Nome, text)
		containsSurName := strings.Contains(m.Pessoa.Sobrenome, text)
		if containsName || containsSurName {
			result = append(result, m)
		}
	}
	return result,nil
}
