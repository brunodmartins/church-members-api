package entity

import (
	"strings"
	"time"
)

type SortByBirthDay []*Membro
type SortByMarriageDay []*Membro
type SortByName []*Membro

func (a SortByName) Len() int      { return len(a) }
func (a SortByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByName) Less(i, j int) bool {
	return strings.Compare(a[i].Pessoa.GetFullName(), a[j].Pessoa.GetFullName()) < 0
}

func (a SortByBirthDay) Len() int      { return len(a) }
func (a SortByBirthDay) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByBirthDay) Less(i, j int) bool {
	firstDate := a[i].Pessoa.DtNascimento
	secondDate := a[j].Pessoa.DtNascimento
	return lessByDay(firstDate, secondDate)
}

func (a SortByMarriageDay) Len() int      { return len(a) }
func (a SortByMarriageDay) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByMarriageDay) Less(i, j int) bool {
	firstDate := a[i].Pessoa.DtCasamento
	secondDate := a[j].Pessoa.DtCasamento
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
