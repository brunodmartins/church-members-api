package entity

import (
	"github.com/bearbin/go-age"
)

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
