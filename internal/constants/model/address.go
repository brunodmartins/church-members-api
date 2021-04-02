package model

import "fmt"

const (
	//HOUSE const
	HOUSE string = "H"
	//APARTMENT const
	APARTMENT string = "A"
)

//Address struct
type Address struct {
	ZipCode  string `json:"zipCode"`
	State    string `json:"state"`
	City     string `json:"city"`
	Address  string `json:"address"`
	District string `json:"district"`
	Number   int    `json:"number"`
	MoreInfo string `json:"moreInfo"`
}

//GetFormatted address
func (address Address) GetFormatted() string {
	return fmt.Sprintf("%s, %d - %s", address.Address, address.Number, address.District)
}
