package jobs

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"time"
)

var genericError = errors.New("generic error")

func BuildBirthDaysMembers() []*domain.Member {
	return []*domain.Member{
		{
			Person: domain.Person{
				BirthDate: time.Now(),
				FirstName: "foo",
				LastName:  "bar",
			},
		},
	}
}

func BuildMarriageMembers() []*domain.Member {
	now := time.Now()
	return []*domain.Member{
		{
			Person: domain.Person{
				MarriageDate: &now,
				FirstName:    "foo",
				LastName:     "bar",
				SpousesName:  "foo2 bar2",
			},
		},
	}
}
