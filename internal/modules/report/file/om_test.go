package file_test

import (
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

func BuildMembers(size int) []*domain.Member {
	var members []*domain.Member
	members = append(members, buildInactiveMember(domain.NewID()))
	for i := 0; i < size; i++ {
		members = append(members, buildMember(domain.NewID()))
	}
	members = append(members, buildInactiveMember(domain.NewID()))
	return members
}

func buildInactiveMember(id string) *domain.Member {
	member := buildMember(id)
	member.Active = false
	member.MembershipEndReason = "End test reason"
	now := time.Now()
	member.MembershipEndDate = &now
	return member
}

func buildMember(id string) *domain.Member {
	now := time.Now()
	return &domain.Member{
		ID:     id,
		Active: true,
		Person: &domain.Person{
			FirstName:    "First Name",
			LastName:     "Last Name",
			BirthDate:    now,
			MarriageDate: &now,
			SpousesName:  "Spouses name",
			Contact: &domain.Contact{
				CellPhoneArea: 99,
				CellPhone:     1234567890,
				PhoneArea:     99,
				Phone:         12345678,
				Email:         "teste@test.com",
			},
			Address: &domain.Address{
				District: "9",
				City:     "Does not sleep",
				State:    "My-State",
				Address:  "X Street",
				Number:   9,
			},
		},
	}
}

func buildChurch() *domain.Church {
	return &domain.Church{
		Name: "teste",
	}
}
