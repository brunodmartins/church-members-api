package entity

import "fmt"

const (
	//CASA const
	CASA string = "C"
	//APARTMENTO const
	APARTMENTO string = "A"
)

//Endereco struct
type Endereco struct {
	Cep        string `json:"cep"`
	UF         string `json:"uf"`
	Cidade     string `json:"cidade"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Numero     int    `json:"numero"`
}

func (this Endereco) GetFormatted() string {
	return fmt.Sprintf("%s, %d - %s", this.Logradouro, this.Numero, this.Bairro)
}
