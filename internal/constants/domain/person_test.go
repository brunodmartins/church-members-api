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
