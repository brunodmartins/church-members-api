package entity

import (
	"fmt"
)

//Contact struct
type Contact struct {
	PhoneArea     int    `json:"phone_area,omitempty" bson:"phone_area"`
	Phone         int    `json:"phone,omitempty"`
	CellPhoneArea int    `json:"cell_phone_area" bson:"cell_phone_area"`
	CellPhone     int    `json:"cell_phone"`
	Email         string `json:"email"`
}

//GetFormattedPhone (99) 99999999
func (c Contact) GetFormattedPhone() string {
	return fmt.Sprintf("(%d) %d", c.PhoneArea, c.Phone)
}

//GetFormattedCellPhone (99) 999999999
func (c Contact) GetFormattedCellPhone() string {
	return fmt.Sprintf("(%d) %d", c.CellPhoneArea, c.CellPhone)
}
