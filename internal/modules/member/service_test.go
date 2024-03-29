package member_test

import (
	"context"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/classification"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	mock_member "github.com/brunodmartins/church-members-api/internal/modules/member/mock"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"testing"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
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
		members, err := service.SearchMembers(BuildContext(), member.OnlyActive())
		assert.Nil(t, err)
		assert.Len(t, members, 2)
	})
	t.Run("Success with post specification", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(2), nil)
		members, err := service.SearchMembers(BuildContext(), member.OnlyActive(), member.OnlyByClassification(classification.CHILDREN))
		assert.Nil(t, err)
		assert.Len(t, members, 2)
	})
	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(nil, genericError)
		_, err := service.SearchMembers(BuildContext(), member.OnlyActive())
		assert.NotNil(t, err)
	})

}

func TestFindMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)

	id := domain.NewID()
	churchMember := buildMember(id)
	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(churchMember, nil)
		found, err := service.GetMember(BuildContext(), id)
		assert.Equal(t, id, found.ID)
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(nil, genericError)
		_, err := service.GetMember(BuildContext(), id)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Invalid ID", func(t *testing.T) {
		_, err := service.GetMember(BuildContext(), "")
		assert.NotNil(t, err)
	})
}

func TestSaveMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)

	t.Run("Success", func(t *testing.T) {
		churchMember := buildMember("")
		repo.EXPECT().Insert(gomock.Any(), gomock.AssignableToTypeOf(churchMember)).DoAndReturn(func(ctx context.Context, member *domain.Member) error {
			member.ID = domain.NewID()
			return nil
		})
		id, err := service.SaveMember(BuildContext(), churchMember)
		assert.Nil(t, err)
		assert.NotEmpty(t, churchMember.ID)
		assert.NotEmpty(t, id)
		assert.NotEmpty(t, churchMember.MembershipStartDate)
		assert.True(t, churchMember.Active)
	})
	t.Run("Fail", func(t *testing.T) {
		churchMember := buildMember("")
		repo.EXPECT().Insert(gomock.Any(), gomock.AssignableToTypeOf(churchMember)).Return(genericError)
		_, err := service.SaveMember(BuildContext(), churchMember)
		assert.NotNil(t, err)
	})
}

func TestChangeStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)
	id := domain.NewID()
	reason := "test"
	date := time.Now()
	churchMember := buildMember(id)
	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(churchMember, nil)
		repo.EXPECT().RetireMembership(gomock.Any(), gomock.AssignableToTypeOf(churchMember)).Return(nil)
		assert.Nil(t, service.RetireMembership(BuildContext(), id, reason, date))
	})
	t.Run("Fail - Retire member", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(churchMember, nil)
		repo.EXPECT().RetireMembership(gomock.Any(), gomock.AssignableToTypeOf(churchMember)).Return(genericError)
		assert.NotNil(t, service.RetireMembership(BuildContext(), id, reason, date))
	})
	t.Run("Fail - Get Member", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(nil, genericError)
		assert.NotNil(t, service.RetireMembership(BuildContext(), id, reason, date))
	})
}
