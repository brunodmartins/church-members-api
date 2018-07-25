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
	DtNascimento     time.Time `json:"dtNascimento"`
	DtCasamento      time.Time `json:"dtCasamento"`
	Naturalidade     string    `json:"naturalidade"`
	CidadeNascimento string    `json:"cidadeNascimento"`
	NomePai          string    `json:"nomePai"`
	NomeMae          string    `json:"nomeMae"`
	NomeConjuge      string    `json:"nomeConjuge"`
	QtdIrmao         int       `json:"qtdIrmao"`
	QtdFilhos        int       `json:"qtdFilhos"`
	Profissao        string    `json:"profissao"`
	Sexo             string    `json:"sexo"`
	EstadoCivil      string    `json:"estadoCivil"`
	Contato          Contato   `json:"contato"`
	Endereco         Endereco  `json:"endereco"`
}
