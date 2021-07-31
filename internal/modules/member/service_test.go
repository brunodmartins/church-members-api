package member_test

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member/mock"
	"testing"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListAllMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)

	repo.EXPECT().FindAll(gomock.AssignableToTypeOf(member.Specification(nil))).Return([]*domain.Member{
		{},
	}, nil).AnyTimes()
	members, _ := service.FindMembers(member.CreateActiveFilter())
	if len(members) == 0 {
		t.Error("No members returned from database")
	}
}

func TestFindMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)

	id := domain.NewID()
	member := &domain.Member{
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
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)

	member := domain.Member{}

	repo.EXPECT().Insert(&member).Return(domain.NewID(), nil).AnyTimes()

	id, err := service.SaveMember(&member)
	if err != nil {
		t.Fail()
	}
	if !domain.IsValidID(id) {
		t.Fail()
	}
}

func TestUpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)
	id := domain.NewID()
	repo.EXPECT().UpdateStatus(id, true).Return(nil)
	repo.EXPECT().GenerateStatusHistory(id, true, "Exited", gomock.Any()).Return(nil)
	err := service.ChangeStatus(id, true, "Exited", time.Now())
	assert.Nil(t, err, "Error not nil")
}

func TestUpdateStatusError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)
	id := domain.NewID()
	repo.EXPECT().UpdateStatus(id, true).Return(errors.New("Error"))
	err := service.ChangeStatus(id, true, "Exited", time.Now())
	assert.NotNil(t, err, "Error not raised")
}
