package entity

import (
	"time"
)

const (
	SOLTEIRO   string = "S"
	CASADO     string = "C"
	DIVORCIADO string = "D"
	VIUVO      string = "V"
)

type Pessoa struct {
	Nome             string    `json:"nome"`
	Sobrenome        string    `json:"sobrenome"`
	DtNascimento     time.Time `json:"dtNascimento" bson:"dtNascimento"`
	DtCasamento      time.Time `json:"dtCasamento,omitempty" bson:"dtCasamento"`
	Naturalidade     string    `json:"naturalidade"`
	CidadeNascimento string    `json:"cidadeNascimento" bson:"cidadeNascimento"`
	NomePai          string    `json:"nomePai" bson:"nomePai"`
	NomeMae          string    `json:"nomeMae" bson:"nomeMae"`
	NomeConjuge      string    `json:"nomeConjuge,omitempty" bson:"nomeConjuge"`
	QtdIrmao         int       `json:"qtdIrmao" bson:"qtdIrmao"`
	QtdFilhos        int       `json:"qtdFilhos" bson:"qtdFilhos"`
	Profissao        string    `json:"profissao,omitempty"`
	Sexo             string    `json:"sexo"`
	Contato          Contato   `json:"contato"`
	Endereco         Endereco  `json:"endereco"`
}

func (pessoa Pessoa) GetFullName() string {
	return pessoa.Nome + " " + pessoa.Sobrenome
}