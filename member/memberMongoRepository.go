package member

import (
	"github.com/BrunoDM2943/church-members-api/entity"
	mgo "gopkg.in/mgo.v2"
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
