package domain

import (
	"fmt"
)

// Contact struct
type Contact struct {
	PhoneArea     int    `json:"phoneArea,omitempty"`
	Phone         int    `json:"phone,omitempty"`
	CellPhoneArea int    `json:"cellPhoneArea"`
	CellPhone     int    `json:"cellPhone"`
	Email         string `json:"email"`
}

// GetFormattedPhone (99) 99999999
func (c Contact) GetFormattedPhone() string {
	if c.Phone == 0 {
		return ""
	}
	return fmt.Sprintf("(%d) %d", c.PhoneArea, c.Phone)
}

// GetFormattedCellPhone (99) 999999999
func (c Contact) GetFormattedCellPhone() string {
	if c.CellPhone == 0 {
		return ""
	}
	return fmt.Sprintf("(%d) %d", c.CellPhoneArea, c.CellPhone)
}
