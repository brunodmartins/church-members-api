package service

import (
	"fmt"
	"testing"

	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/infra/mongo"
	mock_repository "github.com/BrunoDM2943/church-members-api/member/repository/mock"
	"github.com/golang/mock/gomock"
	"gopkg.in/mgo.v2/bson"
)

func TestListAllMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockIMemberRepository(ctrl)
	service := NewMemberService(repo)

	repo.EXPECT().FindAll(gomock.Any()).Return([]*entity.Membro{
		&entity.Membro{},
	}, nil).AnyTimes()
	membros, _ := service.FindMembers(map[string]interface{}{})
	if len(membros) == 0 {
		t.Error("No members returned from database")
	}
}

func TestFindMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockIMemberRepository(ctrl)
	service := NewMemberService(repo)

	id := entity.NewID()
	membro := &entity.Membro{
		ID: id,
	}
	repo.EXPECT().FindByID(id).Return(membro, nil).AnyTimes()
	membroFound, _ := service.FindMembersByID(id)
	if membroFound.ID != id {
		t.Error("Member not found")
	}
}

func TestSaveMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockIMemberRepository(ctrl)
	service := NewMemberService(repo)

	membro := entity.Membro{}

	repo.EXPECT().Insert(&membro).Return(entity.NewID(), nil).AnyTimes()

	id, err := service.SaveMember(&membro)
	if err != nil {
		t.Fail()
	}
	if !entity.IsValidID(id.String()) {
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
	filters.AddFilter("pessoa.sexo", "F")
	regex := bson.RegEx{fmt.Sprintf(".*%s*.", "Bruno"), "i"}
	filters.AddFilter("$or", []bson.M{
		{"pessoa.nome": regex},
		{"pessoa.sobrenome": regex},
	})
	repo.EXPECT().FindAll(gomock.Eq(filters)).Return(nil, nil).AnyTimes()

	service.FindMembers(map[string]interface{}{
		"sexo":   "F",
		"active": true,
		"name":   "Bruno",
	})
}
