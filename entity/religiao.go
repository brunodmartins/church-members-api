package entity

import (
	"time"
)

//Religiao struct
type Religiao struct {
	ReligiaoPais           string    `json:"religiaopais,omitempty"`
	LocalBatismo           string    `json:"localBatismo" bson:"localBatismo"`
	IdadeConheceuEvangelho int       `json:"idadeConheceuEvangelho" bson:"idadeConheceuEvangelho"`
	AceitouJesus           bool      `json:"aceitouJesus" bson:"aceitouJesus"`
	Batizado               bool      `json:"batizado"`
	BatizadoCatolica       bool      `json:"batizadoCatolica" bson:"batizadoCatolica"`
	ConheceDizimo          bool      `json:"conheceDizimo" bson:"conheceDizimo"`
	ConcordaDizimo         bool      `json:"concordaDizimo" bson:"concordaDizimo"`
	Dizimista              bool      `json:"dizimista"`
	DtAceitouJesus         time.Time `json:"dtAceitouJesus" bson:"dtAceitouJesus"`
	DtBatismo              time.Time `json:"dtBatismo" bson:"dtBatismo"`
}
