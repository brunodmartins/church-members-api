package model

import (
	"github.com/magiconair/properties/assert"
	"sort"
	"testing"
	"time"
)

func TestSortByBirth(t *testing.T) {
	firstId := NewID()
	secondId := NewID()
	members := []*Member{
		{
			ID: secondId,
			Person: Person{
				BirthDate: time.Now().AddDate(0, 0, 2),
			},
		},
		{
			ID: firstId,
			Person: Person{
				BirthDate: time.Now().AddDate(0, 0, 1),
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
	members := []*Member{
		{
			ID: secondId,
			Person: Person{
				MarriageDate: time.Now().AddDate(-2, 0, 2),
			},
		},
		{
			ID: firstId,
			Person: Person{
				MarriageDate: time.Now().AddDate(0, 0, 1),
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
				LastName: "Mclane",
			},
		},
		{
			Person: Person{
				FirstName: "John",
				LastName: "Doe",
			},
		},
	}
	sort.Sort(SortByName(members))
	assert.Equal(t, members[0].Person.GetFullName(), "John Doe")
	assert.Equal(t, members[1].Person.GetFullName(), "John Mclane")
}