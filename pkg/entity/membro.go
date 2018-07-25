package entity

type Membro struct {
	ID                     int      `json:"id"`
	AntigaIgreja           string   `json:"antigaIgreja"`
	Ativo                  bool     `json:"bool"`
	FrequentaCultoSexta    bool     `json:"frequentaCultoSexta"`
	FrequentaCultoSabado   bool     `json:"frequentaCultoSabado"`
	FrequentaCultoDomingo  bool     `json:"frequentaCultoDomingo"`
	FrequentaEBD           bool     `json:"frequentaEBD"`
	ImpedimentosFrequencia string   `json:"impedimentosFrequencia"`
	Pessoa                 Pessoa   `json:"pessoa"`
	Religiao               Religiao `json:"religiao"`
}
