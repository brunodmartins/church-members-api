package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPerson_Age(t *testing.T) {
	birthDate, _ := time.Parse("2006-02-01", "2000-01-01")
	currentYear := time.Now().Year()
	person := &Person{BirthDate: birthDate}
	assert.Equal(t, person.Age(), currentYear-birthDate.Year())
}

func TestPerson_GetFullName(t *testing.T) {
	person := &Person{
		FirstName: "John",
		LastName:  "Doe",
	}
	assert.Equal(t, "John Doe", person.GetFullName())
}

func TestPerson_IsMarried(t *testing.T) {
	assert.True(t, Person{MaritalStatus: "MARRIED"}.IsMarried())
	assert.False(t, Person{MaritalStatus: "SINGLE"}.IsMarried())
}

func TestPerson_GetCoupleName(t *testing.T) {
	assert.Equal(t, "John Doe & Mary", Person{
		FirstName:   "John",
		LastName:    "Doe",
		SpousesName: "Mary",
	}.GetCoupleName())
}
