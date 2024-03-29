package domain

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSortByBirth(t *testing.T) {
	firstId := NewID()
	secondId := NewID()
	thirdId := NewID()
	fourthId := NewID()
	baseTestTime, _ := time.Parse(time.RFC3339, "2020-01-01")
	members := []*Member{
		{
			ID: thirdId,
			Person: &Person{
				BirthDate: baseTestTime.AddDate(0, 1, 10),
			},
		},
		{
			ID: fourthId,
			Person: &Person{
				BirthDate: baseTestTime.AddDate(0, 2, 2),
			},
		},
		{
			ID: secondId,
			Person: &Person{
				BirthDate: baseTestTime.AddDate(0, 0, 2),
			},
		},
		{
			ID: firstId,
			Person: &Person{
				BirthDate: baseTestTime.AddDate(0, 0, 1),
			},
		},
	}
	sort.Sort(SortByBirthDay(members))
	assert.Equal(t, members[0].ID, firstId)
	assert.Equal(t, members[1].ID, secondId)
	assert.Equal(t, members[2].ID, thirdId)
	assert.Equal(t, members[3].ID, fourthId)
}

func TestSortByMarriage(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	firstId := NewID()
	secondId := NewID()
	birthDateOne := now.AddDate(-2, 0, 2)
	birthDateTwo := now.AddDate(0, 0, 1)
	members := []*Member{
		{
			ID: secondId,
			Person: &Person{
				MarriageDate: &birthDateOne,
			},
		},
		{
			ID: firstId,
			Person: &Person{
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
			Person: &Person{
				FirstName: "John",
				LastName:  "Mclane",
			},
		},
		{
			Person: &Person{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
	}
	sort.Sort(SortByName(members))
	assert.Equal(t, members[0].Person.GetFullName(), "John Doe")
	assert.Equal(t, members[1].Person.GetFullName(), "John Mclane")
}

func TestLessByDay(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	assert.True(t, lessByDay(now, now.AddDate(0, 1, 0)))
	assert.False(t, lessByDay(now.AddDate(0, 1, 0), now))
	assert.True(t, lessByDay(now, now.AddDate(0, 0, 1)))
	assert.False(t, lessByDay(now.AddDate(0, 0, 1), now))
}
