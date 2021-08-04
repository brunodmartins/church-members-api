package domain

import (
	"sort"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestSortByBirth(t *testing.T) {
	firstId := NewID()
	secondId := NewID()
	birthDateOne := time.Now().AddDate(0, 0, 2)
	birthDateTwo := time.Now().AddDate(0, 0, 1)
	members := []*Member{
		{
			ID: secondId,
			Person: Person{
				BirthDate: birthDateOne,
			},
		},
		{
			ID: firstId,
			Person: Person{
				BirthDate: birthDateTwo,
			},
		},
	}
	sort.Sort(SortByBirthDay(members))
	assert.Equal(t, members[0].ID, firstId)
	assert.Equal(t, members[1].ID, secondId)
}

func TestSortByMarriage(t *testing.T) {
	firstId := NewID()
	secondId := NewID()
	birthDateOne := time.Now().AddDate(-2, 0, 2)
	birthDateTwo := time.Now().AddDate(0, 0, 1)
	members := []*Member{
		{
			ID: secondId,
			Person: Person{
				MarriageDate: &birthDateOne,
			},
		},
		{
			ID: firstId,
			Person: Person{
				MarriageDate: &birthDateTwo,
			},
		},
	}
	sort.Sort(SortByMarriageDay(members))
	assert.Equal(t, members[0].ID, firstId)
	assert.Equal(t, members[1].ID, secondId)
}

func TestSortByName(t *testing.T) {
	members := []*Member{
		{
			Person: Person{
				FirstName: "John",
				LastName:  "Mclane",
			},
		},
		{
			Person: Person{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
	}
	sort.Sort(SortByName(members))
	assert.Equal(t, members[0].Person.GetFullName(), "John Doe")
	assert.Equal(t, members[1].Person.GetFullName(), "John Mclane")
}
