package entity

const (
	//CASA const
	CASA string = "C"
	//APARTMENTO const
	APARTMENTO string = "A"
)

//Endereco struct
type Endereco struct {
	Cep            string `json:"cep"`
	UF             string `json:"uf"`
	Cidade         string `json:"cidade"`
	Logradouro     string `json:"logradouro"`
	Bairro         string `json:"bairro"`
	TipoResidencia string `json:"tipoResidencia"`
	Numero         int    `json:"numero"`
}
