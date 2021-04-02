package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	mock_repository "github.com/BrunoDM2943/church-members-api/internal/repository/mock"
	"github.com/BrunoDM2943/church-members-api/internal/storage/mongo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestListAllMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockIMemberRepository(ctrl)
	service := NewMemberService(repo)

	repo.EXPECT().FindAll(gomock.Any()).Return([]*model.Member{
		&model.Member{},
	}, nil).AnyTimes()
	members, _ := service.FindMembers(map[string]interface{}{})
	if len(members) == 0 {
		t.Error("No members returned from database")
	}
}

func TestFindMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockIMemberRepository(ctrl)
	service := NewMemberService(repo)

	id := model.NewID()
	member := &model.Member{
		ID: id,
	}
	repo.EXPECT().FindByID(id).Return(member, nil).AnyTimes()
	memberFound, _ := service.FindMembersByID(id)
	if memberFound.ID != id {
		t.Error("Member not found")
	}
}

func TestSaveMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockIMemberRepository(ctrl)
	service := NewMemberService(repo)

	member := model.Member{}

	repo.EXPECT().Insert(&member).Return(model.NewID(), nil).AnyTimes()

	id, err := service.SaveMember(&member)
	if err != nil {
		t.Fail()
	}
	if !model.IsValidID(id.String()) {
		t.Fail()
	}
}

func TestFindMembersWithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockIMemberRepository(ctrl)
	service := NewMemberService(repo)
	filters := mongo.QueryFilters{}
	filters.AddFilter("active", true)
	filters.AddFilter("person.gender", "F")
	regex := bson.RegEx{fmt.Sprintf(".*%s*.", "Bruno"), "i"}
	filters.AddFilter("$or", []bson.M{
		{"person.firstName": regex},
		{"person.lastName": regex},
	})
	repo.EXPECT().FindAll(gomock.Eq(filters)).Return(nil, nil).AnyTimes()

	service.FindMembers(map[string]interface{}{
		"gender": "F",
		"active": true,
		"name":   "Bruno",
	})
}

func TestUpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockIMemberRepository(ctrl)
	service := NewMemberService(repo)
	id := model.NewID()
	repo.EXPECT().UpdateStatus(id, true).Return(nil)
	repo.EXPECT().GenerateStatusHistory(id, true, "Exited", gomock.Any()).Return(nil)
	err := service.ChangeStatus(id, true, "Exited")
	assert.Nil(t, err, "Error not nil")
}

func TestUpdateStatusError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockIMemberRepository(ctrl)
	service := NewMemberService(repo)
	id := model.NewID()
	repo.EXPECT().UpdateStatus(id, true).Return(errors.New("Error"))
	repo.EXPECT().GenerateStatusHistory(id, true, "Exited", gomock.Any()).Return(nil)
	err := service.ChangeStatus(id, true, "Exited")
	assert.NotNil(t, err, "Error not raised")
}
