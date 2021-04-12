package repository

import (
	"fmt"
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type mongoRepository struct {
	col        *mgo.Collection
	colHistory *mgo.Collection
}

func NewMongoRepository(session *mgo.Session) MemberRepository {
	return &mongoRepository{
		col:        session.DB("disciples").C("member"),
		colHistory: session.DB("disciples").C("member_history"),
	}
}

func (repo *mongoRepository) FindAll(filters QueryFilters) ([]*model.Member, error) {
	var result []*model.Member

	if filters.GetFilter("name") != nil {
		name := filters.GetFilter("name").(string)
		regex := bson.RegEx{fmt.Sprintf(".*%s*.", name), "i"}
		filters.AddFilter("$or", []bson.M{
			{"person.firstName": regex},
			{"person.lastName": regex},
		})
	}

	err := repo.col.Find(filters).Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *mongoRepository) FindByID(id model.ID) (*model.Member, error) {
	var result *model.Member
	err := repo.col.FindId(bson.ObjectIdHex(id.String())).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			return nil, MemberNotFound
		}
		return nil, err
	}
	return result, nil
}

func (repo *mongoRepository) Insert(member *model.Member) (model.ID, error) {
	member.ID = model.NewID()
	return member.ID, repo.col.Insert(member)
}

func (repo *mongoRepository) UpdateStatus(ID model.ID, status bool) error {
	return repo.col.UpdateId(bson.ObjectIdHex(ID.String()), bson.M{
		"$set": bson.M{
			"active": status,
		}})
}

func (repo *mongoRepository) GenerateStatusHistory(id model.ID, status bool, reason string, date time.Time) error {
	return repo.colHistory.Insert(bson.M{
		"member_id":    id,
		"status":       status,
		"reason":       reason,
		"changed_date": date,
	})
}

func (repo *mongoRepository) FindMembersActive() ([]*model.Member, error) {
	var result []*model.Member
	filters := QueryFilters{}
	filters.AddFilter("active", true)
	err := repo.col.Find(filters).Sort("person.firstName", "person.lastName").Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (repo *mongoRepository) FindMembersActiveAndMarried() ([]*model.Member, error) {
	var result []*model.Member
	filters := QueryFilters{}
	filters.AddFilter("active", true)
	filters.AddFilter("person.marriageDate", bson.M{"$exists": true})
	err := repo.col.Find(filters).Select(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}