package jobs

import (
	"errors"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
)

var genericError = errors.New("generic error")

func BuildBirthDaysMembers(date time.Time) []*domain.Member {
	return []*domain.Member{
		{
			Person: domain.Person{
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
			Person: domain.Person{
				MarriageDate: date,
				FirstName:    "foo",
				LastName:     "bar",
				SpousesName:  "foo2 bar2",
			},
		},
	}
}
