package entity

import "fmt"

const (
	//HOUSE const
	HOUSE string = "H"
	//APARTMENT const
	APARTMENT string = "A"
)

//Address struct
type Address struct {
	ZipCode  string `json:"zip_code"`
	State    string `json:"state"`
	City     string `json:"city"`
	Address  string `json:"address"`
	District string `json:"district"`
	Number   int    `json:"number"`
}

//GetFormatted address
func (address Address) GetFormatted() string {
	return fmt.Sprintf("%s, %d - %s", address.Address, address.Number, address.District)
}
