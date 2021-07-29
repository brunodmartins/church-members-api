package member

import (
	"errors"
	"testing"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/repository"

	"github.com/BrunoDM2943/church-members-api/internal/constants/entity"
	mock_repository "github.com/BrunoDM2943/church-members-api/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListAllMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	service := NewMemberService(repo)

	repo.EXPECT().FindAll(gomock.Any()).Return([]*entity.Member{
		{},
	}, nil).AnyTimes()
	members, _ := service.FindMembers(map[string]interface{}{})
	if len(members) == 0 {
		t.Error("No members returned from database")
	}
}

func TestFindMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	service := NewMemberService(repo)

	id := entity.NewID()
	member := &entity.Member{
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
	repo := mock_repository.NewMockMemberRepository(ctrl)
	service := NewMemberService(repo)

	member := entity.Member{}

	repo.EXPECT().Insert(&member).Return(entity.NewID(), nil).AnyTimes()

	id, err := service.SaveMember(&member)
	if err != nil {
		t.Fail()
	}
	if !entity.IsValidID(id) {
		t.Fail()
	}
}

func TestFindMembersWithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	service := NewMemberService(repo)
	filters := repository.QueryFilters{}
	filters.AddFilter("active", true)
	filters.AddFilter("person.gender", "F")
	filters.AddFilter("name", "Bruno")
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
	repo := mock_repository.NewMockMemberRepository(ctrl)
	service := NewMemberService(repo)
	id := entity.NewID()
	repo.EXPECT().UpdateStatus(id, true).Return(nil)
	repo.EXPECT().GenerateStatusHistory(id, true, "Exited", gomock.Any()).Return(nil)
	err := service.ChangeStatus(id, true, "Exited", time.Now())
	assert.Nil(t, err, "Error not nil")
}

func TestUpdateStatusError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	service := NewMemberService(repo)
	id := entity.NewID()
	repo.EXPECT().UpdateStatus(id, true).Return(errors.New("Error"))
	err := service.ChangeStatus(id, true, "Exited", time.Now())
	assert.NotNil(t, err, "Error not raised")
}
