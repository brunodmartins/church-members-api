package jobs

import (
	"errors"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

var genericError = errors.New("generic error")

func BuildBirthDaysMembers(date time.Time) []*domain.Member {
	return []*domain.Member{
		{
			Person: &domain.Person{
				BirthDate: date,
				FirstName: "foo",
				LastName:  "bar",
			},
		},
	}
}

func BuildMarriageMembers(date *time.Time) []*domain.Member {

	return []*domain.Member{
		{
			Person: &domain.Person{
				MarriageDate: date,
				FirstName:    "foo",
				LastName:     "bar",
				SpousesName:  "foo2 bar2",
			},
		},
	}
}

func BuildUsers() []*domain.User {
	return []*domain.User{
		{
			Email: "test@test.com",
			Phone: "123",
		},
		{
			Email: "test@test.com",
			Phone: "123",
		},
	}
}

func buildChurchs() []*domain.Church {
	return []*domain.Church{
		{
			ID:       "1",
			Language: "pt-BR",
		}, {
			ID:       "2",
			Language: "en",
		},
	}
}
