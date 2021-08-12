package domain

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormattedContact(t *testing.T) {
	t.Run("With Values", func(t *testing.T) {
		contact := Contact{
			CellPhone:     953200587,
			CellPhoneArea: 11,
			Phone:         29435002,
			PhoneArea:     11,
		}
		assert.Equal(t, "(11) 953200587", contact.GetFormattedCellPhone())
		assert.Equal(t, "(11) 29435002", contact.GetFormattedPhone())
	})
	t.Run("Without Values", func(t *testing.T) {
		contact := Contact{}
		assert.Empty(t, contact.GetFormattedCellPhone())
		assert.Empty(t, contact.GetFormattedPhone())
	})
}

func TestClassification(t *testing.T) {
	t.Run("Children", func(t *testing.T) {
		assert.Equal(t, enum.CHILDREN, Member{
			Person: Person{
				BirthDate: time.Now(),
			},
		}.Classification())
	})
	t.Run("Teen", func(t *testing.T) {
		assert.Equal(t, enum.TEEN, Member{
			Person: Person{
				BirthDate: time.Now().AddDate(-17, 0, 0),
			},
		}.Classification())
	})
	t.Run("Young", func(t *testing.T) {
		assert.Equal(t, enum.YOUNG, Member{
			Person: Person{
				BirthDate: time.Now().AddDate(-29, 0, 0),
			},
		}.Classification())
	})
	t.Run("Adult Single", func(t *testing.T) {
		assert.Equal(t, enum.ADULT, Member{
			Person: Person{
				BirthDate: time.Now().AddDate(-33, 0, 0),
			},
		}.Classification())
	})
	t.Run("Adult Married", func(t *testing.T) {
		now := time.Now()
		assert.Equal(t, enum.ADULT, Member{
			Person: Person{
				BirthDate:    time.Now().AddDate(-25, 0, 0),
				MarriageDate: &now,
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
	assert.Equal(t, "Rua xicas, 2 - Parque feliz", address.String())
}

func TestGetFullName(t *testing.T) {
	assert.Equal(t, Person{
		FirstName: "John",
		LastName:  "Doe",
	}.GetFullName(), "John Doe")
}