package entity

import (
	"testing"
)

func TestFormattedContact(t *testing.T) {
	c := Contato{
		Celular:    953200587,
		DDDCelular: 11,
	}
	if "(11) 953200587" != c.GetFormattedCellPhone() {
		t.Fail()
	}
}
