package member

import (
	"fmt"
	"github.com/BrunoDM2943/church-members-api/entity"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MemberRepository struct {
	col *mgo.Collection
}

func NewMemberRepository(session *mgo.Session) *MemberRepository {
	return &MemberRepository{
		col: session.DB("disciples").C("Membro"),
	}
}

func (repo *MemberRepository) FindAll() ([]*entity.Membro, error) {
	var result []*entity.Membro
	err := repo.col.Find(nil).Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *MemberRepository) FindByID(id entity.ID) (*entity.Membro, error) {
	var result *entity.Membro
	err := repo.col.FindId(bson.ObjectIdHex(id.String())).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			return nil, MemberNotFound
		}
		return nil, err
	}
	return result, nil
}

func (repo *MemberRepository) Insert(membro *entity.Membro) (entity.ID, error) {
	membro.ID = entity.NewID()
	return membro.ID, repo.col.Insert(membro)
}

func (repo *MemberRepository) Search(text string) ([]*entity.Membro, error) {
	var result []*entity.Membro
	regex := bson.RegEx{fmt.Sprintf(".*%s*.",text), "i"}
	err := repo.col.Find(
		bson.M{
			"$or": []bson.M{
				{"pessoa.nome": regex},
				{"pessoa.sobrenome": regex},
			},
		},
		).Select(bson.M{}).All(&result)
	return result, err
}
