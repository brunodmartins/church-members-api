package entity

import (
	"fmt"
)

//Contato struct
type Contato struct {
	DDDTelefone int    `json:"dddTelefone,omitempty" bson:"dddTelefone"`
	Telefone    int    `json:"telefone,omitempty"`
	DDDCelular  int    `json:"dddCelular" bson:"dddCelular"`
	Celular     int    `json:"celular"`
	Email       string `json:"email"`
}

//GetFormattedPhone (99) 99999999
func (c Contato) GetFormattedPhone() string {
	return fmt.Sprintf("(%d) %d", c.DDDTelefone, c.Telefone)
}

//GetFormattedCellPhone (99) 999999999
func (c Contato) GetFormattedCellPhone() string {
	return fmt.Sprintf("(%d) %d", c.DDDCelular, c.Celular)
}
