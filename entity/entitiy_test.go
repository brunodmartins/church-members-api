package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormattedContact(t *testing.T) {
	c := Contact{
		CellPhone:     953200587,
		CellPhoneArea: 11,
		Phone:         29435002,
		PhoneArea:     11,
	}
	if "(11) 953200587" != c.GetFormattedCellPhone() {
		t.Fail()
	}

	if "(11) 29435002" != c.GetFormattedPhone() {
		t.Fail()
	}
}

func TestClassificacao(t *testing.T) {
	t.Run("Crianca", func(t *testing.T) {
		assert.Equal(t, "Criança", Member{
			Person: Person{
				BirthDate: time.Now(),
			},
		}.Classification())
	})
	t.Run("Adolescente", func(t *testing.T) {
		assert.Equal(t, "Adolescente", Member{
			Person: Person{
				BirthDate: time.Now().AddDate(-17, 0, 0),
			},
		}.Classification())
	})
	t.Run("Jovem", func(t *testing.T) {
		assert.Equal(t, "Jovem", Member{
			Person: Person{
				BirthDate: time.Now().AddDate(-29, 0, 0),
			},
		}.Classification())
	})
	t.Run("Adulto Solteiro", func(t *testing.T) {
		assert.Equal(t, "Adulto", Member{
			Person: Person{
				BirthDate: time.Now().AddDate(-33, 0, 0),
			},
		}.Classification())
	})
	t.Run("Adulto Casado", func(t *testing.T) {
		assert.Equal(t, "Adulto", Member{
			Person: Person{
				BirthDate:    time.Now().AddDate(-25, 0, 0),
				MarriageDate: time.Now(),
			},
		}.Classification())
	})
}

func TestFormattedAddress(t *testing.T) {
	address := Address{
		Address:  "Rua xicas",
		District: "Parque feliz",
		Number:   2,
	}
	assert.Equal(t, "Rua xicas, 2 - Parque feliz", address.GetFormatted())
}
