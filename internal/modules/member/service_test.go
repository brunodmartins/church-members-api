package member_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/enum/classification"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	mock_member "github.com/brunodmartins/church-members-api/internal/modules/member/mock"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"go.uber.org/mock/gomock"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
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

func TestMemberService_UpdateContact(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.TODO()
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)
	t.Run("Successfully update the contact information", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		contact := domain.Contact{
			Email: "test@test.com",
		}
		repo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(churchMember, nil)
		repo.EXPECT().UpdateContact(gomock.Eq(ctx), memberContactMatcher{contact: contact}).Return(nil)
		assert.Nil(t, service.UpdateContact(ctx, id, contact))
	})
	t.Run("Fails to update the contact information due to update error", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		contact := domain.Contact{
			Email: "test@test.com",
		}
		repo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(churchMember, nil)
		repo.EXPECT().UpdateContact(gomock.Eq(ctx), memberContactMatcher{contact: contact}).Return(genericError)
		assert.NotNil(t, service.UpdateContact(ctx, id, contact))
	})
	t.Run("Fails to update the contact information due to find error", func(t *testing.T) {
		id := domain.NewID()
		contact := domain.Contact{
			Email: "test@test.com",
		}
		repo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(nil, genericError)
		assert.NotNil(t, service.UpdateContact(ctx, id, contact))
	})
}

func TestMemberService_UpdateAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.TODO()
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)
	id := domain.NewID()
	churchMember := buildMember(id)
	address := domain.Address{
		ZipCode:  "123456-789",
		State:    "Sao Paulo",
		City:     "Sao Paulo",
		Address:  "Test",
		District: "Testing",
		Number:   123456,
		MoreInfo: "1999",
	}
	t.Run("Successfully update the address information", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(churchMember, nil)
		repo.EXPECT().UpdateAddress(gomock.Eq(ctx), memberAddressMatcher{address: address}).Return(nil)
		assert.NoError(t, service.UpdateAddress(ctx, id, address))
	})
	t.Run("Fails to update the address information due to update error", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(churchMember, nil)
		repo.EXPECT().UpdateAddress(gomock.Eq(ctx), memberAddressMatcher{address: address}).Return(genericError)
		assert.Error(t, service.UpdateAddress(ctx, id, address))
	})
	t.Run("Fails to update the address information due to find error", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(nil, genericError)
		assert.Error(t, service.UpdateAddress(ctx, id, address))
	})
}

func TestMemberService_UpdatePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.TODO()
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)
	id := domain.NewID()
	churchMember := buildMember(id)
	person := domain.Person{
		FirstName:        "New name",
		LastName:         "New last name",
		BirthDate:        churchMember.Person.BirthDate,
		MarriageDate:     nil,
		SpousesName:      "",
		MaritalStatus:    "SINGLE",
		ChildrenQuantity: churchMember.Person.ChildrenQuantity,
	}
	t.Run("Successfully update the person information", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(churchMember, nil)
		repo.EXPECT().UpdatePerson(gomock.Eq(ctx), memberPersonMatcher{person: person}).Return(nil)
		assert.NoError(t, service.UpdatePerson(ctx, id, person))
	})
	t.Run("Fails to update the person information due to update error", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(churchMember, nil)
		repo.EXPECT().UpdatePerson(gomock.Eq(ctx), memberPersonMatcher{person: person}).Return(genericError)
		assert.Error(t, service.UpdatePerson(ctx, id, person))
	})
	t.Run("Fails to update the person information due to find error", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Eq(id)).Return(nil, genericError)
		assert.Error(t, service.UpdatePerson(ctx, id, person))
	})
}

func TestMemberService_GetLastBirthAnniversaries(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.TODO()
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)

	t.Run("Successfully get last birthday members", func(t *testing.T) {
		members := BuildMembers(2)
		repo.EXPECT().FindAll(gomock.Eq(ctx), gomock.Any()).Return(members, nil)
		result, err := service.GetLastBirthAnniversaries(ctx)
		assert.NoError(t, err)
		assert.Equal(t, members, result)
	})

	t.Run("Fails to get last birthday members", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Eq(ctx), gomock.Any()).Return(nil, genericError)
		result, err := service.GetLastBirthAnniversaries(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Returns sorted birthday members", func(t *testing.T) {
		member1 := buildMember("")
		member1.Person.BirthDate = time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)
		member2 := buildMember("")
		member2.Person.BirthDate = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		members := []*domain.Member{member1, member2}

		repo.EXPECT().FindAll(gomock.Eq(ctx), gomock.Any()).Return(members, nil)
		result, err := service.GetLastBirthAnniversaries(ctx)

		assert.NoError(t, err)
		assert.Equal(t, member2.Person.BirthDate, result[0].Person.BirthDate)
		assert.Equal(t, member1.Person.BirthDate, result[1].Person.BirthDate)
	})
}

func TestMemberService_GetLastMarriageAnniversaries(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.TODO()
	defer ctrl.Finish()
	repo := mock_member.NewMockRepository(ctrl)
	service := member.NewMemberService(repo)

	t.Run("Successfully get last marriage members", func(t *testing.T) {
		members := BuildMembers(2)
		repo.EXPECT().FindAll(gomock.Eq(ctx), gomock.Any()).Return(members, nil)
		result, err := service.GetLastMarriageAnniversaries(ctx)
		assert.NoError(t, err)
		assert.Equal(t, members, result)
	})

	t.Run("Fails to get last marriage members", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Eq(ctx), gomock.Any()).Return(nil, genericError)
		result, err := service.GetLastMarriageAnniversaries(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Returns sorted marriage members", func(t *testing.T) {
		marriageDate1 := time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)
		marriageDate2 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		member1 := buildMember("")
		member1.Person.MarriageDate = &marriageDate1
		member2 := buildMember("")
		member2.Person.MarriageDate = &marriageDate2
		members := []*domain.Member{member1, member2}

		repo.EXPECT().FindAll(gomock.Eq(ctx), gomock.Any()).Return(members, nil)
		result, err := service.GetLastMarriageAnniversaries(ctx)

		assert.NoError(t, err)
		assert.Equal(t, marriageDate2, *result[0].Person.MarriageDate)
		assert.Equal(t, marriageDate1, *result[1].Person.MarriageDate)
	})
}

type memberContactMatcher struct {
	contact domain.Contact
}

func (expected memberContactMatcher) Matches(received any) bool {
	return *received.(*domain.Member).Person.Contact == expected.contact
}

func (expected memberContactMatcher) String() string {
	return fmt.Sprintf("Expetected %v", expected.contact)
}

type memberAddressMatcher struct {
	address domain.Address
}

func (expected memberAddressMatcher) Matches(received any) bool {
	return *received.(*domain.Member).Person.Address == expected.address
}

func (expected memberAddressMatcher) String() string {
	return fmt.Sprintf("Expetected %v", expected.address)
}

type memberPersonMatcher struct {
	person domain.Person
}

func (expected memberPersonMatcher) Matches(received any) bool {
	receivedMember := received.(*domain.Member)

	if receivedMember.Person.FirstName != expected.person.FirstName {
		return false
	}
	if receivedMember.Person.LastName != expected.person.LastName {
		return false
	}
	if receivedMember.Person.BirthDate != expected.person.BirthDate {
		return false
	}
	if receivedMember.Person.MarriageDate != expected.person.MarriageDate {
		return false
	}
	if receivedMember.Person.ChildrenQuantity != expected.person.ChildrenQuantity {
		return false
	}
	if receivedMember.Person.SpousesName != expected.person.SpousesName {
		return false
	}
	if receivedMember.Person.MaritalStatus != expected.person.MaritalStatus {
		return false
	}

	return true
}

func (expected memberPersonMatcher) String() string {
	return fmt.Sprintf("Expetected %v", expected.person)
}
