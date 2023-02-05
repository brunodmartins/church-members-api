package report_test

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"time"
)

var (
	genericError = errors.New("generic error")
)

func BuildMembers(size int) []*domain.Member {
	var members []*domain.Member
	for i := 0; i < size; i++ {
		members = append(members, buildMember(domain.NewID()))
	}
	return members
}

func buildMember(id string) *domain.Member {
	now := time.Now()
	return &domain.Member{
		ID: id,
		Person: domain.Person{
			FirstName:    "First Name",
			LastName:     "Last Name",
			BirthDate:    now,
			MarriageDate: &now,
			SpousesName:  "Spouses name",
			Contact: domain.Contact{
				CellPhoneArea: 99,
				CellPhone:     1234567890,
				PhoneArea:     99,
				Phone:         12345678,
				Email:         "teste@test.com",
			},
			Address: domain.Address{
				District: "9",
				City:     "Does not sleep",
				State:    "My-State",
				Address:  "X Street",
				Number:   9,
			},
		},
	}
}
