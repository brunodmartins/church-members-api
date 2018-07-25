package entity

import (
	"time"
)

//Religiao struct
type Religiao struct {
	ReligiaoPais           string    `json:"religiaopais"`
	LocalBatismo           string    `json:"localBatismo"`
	IdadeConheceuEvangelho int       `json:"idadeConheceuEvangelho"`
	AceitouJesus           bool      `json:"aceitouJesus"`
	Batizado               bool      `json:"batizado"`
	BatizadoCatolica       bool      `json:"batizadoCatolica"`
	ConheceDizimo          bool      `json:"conheceDizimo"`
	ConcordaDizimo         bool      `json:"concordaDizimo"`
	Dizimista              bool      `json:"dizimista"`
	DtAceitouJesus         time.Time `json:"dtAceitouJesus"`
	DtBatismo              time.Time `json:"dtBatismo"`
}
