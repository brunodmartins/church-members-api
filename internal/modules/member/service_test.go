package member_test

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member/mock"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
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
	spec := wrapper.QuerySpecification(nil)
	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(2), nil)
		members, err := service.SearchMembers(context.Background(), member.OnlyActive())
		assert.Nil(t, err)
		assert.Len(t, members, 2)
	})
	t.Run("Success with post specification", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(2), nil)
		members, err := service.SearchMembers(context.Background(), member.OnlyActive(), member.OnlyLegalMembers())
		assert.Nil(t, err)
		assert.Len(t, members, 2)
	})
	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(nil, genericError)
		_, err := service.SearchMembers(context.Background(), member.OnlyActive())
		assert.NotNil(t, err)
	})

}

func TestFindMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)

	id := domain.NewID()
	member := buildMember(id)
	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(member, nil)
		found, err := service.GetMember(context.Background(), id)
		assert.Equal(t, id, found.ID)
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(nil, genericError)
		_, err := service.GetMember(context.Background(), id)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Invalid ID", func(t *testing.T) {
		_, err := service.GetMember(context.Background(), "")
		assert.NotNil(t, err)
	})
}

func TestSaveMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)

	t.Run("Success", func(t *testing.T) {
		member := buildMember("")
		repo.EXPECT().Insert(gomock.Any(), gomock.AssignableToTypeOf(member)).DoAndReturn(func(ctx context.Context, member *domain.Member) error {
			member.ID = domain.NewID()
			return nil
		})
		id, err := service.SaveMember(context.Background(), member)
		assert.Nil(t, err)
		assert.NotEmpty(t, member.ID)
		assert.NotEmpty(t, id)
		assert.True(t, member.Active)
	})
	t.Run("Fail", func(t *testing.T) {
		member := buildMember("")
		repo.EXPECT().Insert(gomock.Any(), gomock.AssignableToTypeOf(member)).Return(genericError)
		_, err := service.SaveMember(context.Background(), member)
		assert.NotNil(t, err)
	})
}

func TestChangeStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)
	id := domain.NewID()
	status := true
	reason := "test"
	date := time.Now()
	member := buildMember(id)
	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(member, nil)
		repo.EXPECT().UpdateStatus(gomock.Any(), gomock.AssignableToTypeOf(member)).Return(nil)
		repo.EXPECT().GenerateStatusHistory(id, status, reason, date).Return(nil)
		assert.Nil(t, service.ChangeStatus(context.Background(), id, status, reason, date))
	})
	t.Run("Fail - Status History", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(member, nil)
		repo.EXPECT().UpdateStatus(gomock.Any(), gomock.AssignableToTypeOf(member)).Return(nil)
		repo.EXPECT().GenerateStatusHistory(id, status, reason, date).Return(genericError)
		assert.NotNil(t, service.ChangeStatus(context.Background(), id, status, reason, date))
	})
	t.Run("Fail - Update Status", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(member, nil)
		repo.EXPECT().UpdateStatus(gomock.Any(), gomock.AssignableToTypeOf(member)).Return(genericError)
		assert.NotNil(t, service.ChangeStatus(context.Background(), id, status, reason, date))
	})
	t.Run("Fail - Get Member", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(nil, genericError)
		assert.NotNil(t, service.ChangeStatus(context.Background(), id, status, reason, date))
	})
}
