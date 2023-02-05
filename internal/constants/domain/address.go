package domain

import "fmt"

// Address struct
type Address struct {
	ZipCode  string `json:"zipCode"`
	State    string `json:"state"`
	City     string `json:"city"`
	Address  string `json:"address"`
	District string `json:"district"`
	Number   int    `json:"number"`
	MoreInfo string `json:"moreInfo"`
}

// String address
func (address Address) String() string {
	return fmt.Sprintf("%s, %d - %s", address.Address, address.Number, address.District)
}
