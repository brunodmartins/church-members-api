package entity

import (
	"strings"
	"time"
)

//SortByBirthDay a list of members
type SortByBirthDay []*Member

//SortByMarriageDay a list of members
type SortByMarriageDay []*Member

//SortByName a list of members
type SortByName []*Member

func (a SortByName) Len() int      { return len(a) }
func (a SortByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByName) Less(i, j int) bool {
	return strings.Compare(a[i].Person.GetFullName(), a[j].Person.GetFullName()) < 0
}

func (a SortByBirthDay) Len() int      { return len(a) }
func (a SortByBirthDay) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByBirthDay) Less(i, j int) bool {
	firstDate := a[i].Person.BirthDate
	secondDate := a[j].Person.BirthDate
	return lessByDay(firstDate, secondDate)
}

func (a SortByMarriageDay) Len() int      { return len(a) }
func (a SortByMarriageDay) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByMarriageDay) Less(i, j int) bool {
	firstDate := a[i].Person.MarriageDate
	secondDate := a[j].Person.MarriageDate
	return lessByDay(firstDate, secondDate)
}

func lessByDay(firstDate, secondDate time.Time) bool {
	if firstDate.Month() < secondDate.Month() {
		return true
	} else if firstDate.Month() > secondDate.Month() {
		return false
	} else {
		return firstDate.Day() < secondDate.Day()
	}
}
