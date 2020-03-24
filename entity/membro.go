package entity

import (
	"time"

	"github.com/bearbin/go-age"
)

type SortByBirthDay []*Membro

func (a SortByBirthDay) Len() int      { return len(a) }
func (a SortByBirthDay) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByBirthDay) Less(i, j int) bool {
	firstDate := a[i].Pessoa.DtNascimento
	secondDate := a[j].Pessoa.DtNascimento
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

type SortByMarriageDay []*Membro

func (a SortByMarriageDay) Len() int      { return len(a) }
func (a SortByMarriageDay) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByMarriageDay) Less(i, j int) bool {
	firstDate := a[i].Pessoa.DtCasamento
	secondDate := a[j].Pessoa.DtCasamento
	return lessByDay(firstDate, secondDate)
}

type Membro struct {
	ID                     ID       `json:"id" bson:"_id,omitempty"`
	AntigaIgreja           string   `json:"antigaIgreja,omitempty"`
	FrequentaCultoSexta    bool     `json:"frequentaCultoSexta" bson:"frequentaCultoSexta"`
	FrequentaCultoSabado   bool     `json:"frequentaCultoSabado" bson:"frequentaCultoSabado"`
	FrequentaCultoDomingo  bool     `json:"frequentaCultoDomingo" bson:"frequentaCultoDomingo"`
	FrequentaEBD           bool     `json:"frequentaEBD" bson:"frequentaEBD"`
	ImpedimentosFrequencia string   `json:"impedimentosFrequencia,omitempty" bson:"impedimentosFrequencia"`
	Pessoa                 Pessoa   `json:"pessoa"`
	Religiao               Religiao `json:"religiao"`
	Active                 bool     `json:"active,omitempty" bson:"active"`
}

func (this Membro) Classificacao() string {
	age := age.Age(this.Pessoa.DtNascimento)
	if age < 15 {
		return "Criança"
	} else if age < 18 {
		return "Adolescente"
	} else if age < 30 && this.Pessoa.DtCasamento.IsZero() {
		return "Jovem"
	} else {
		return "Adulto"
	}
}
