package domain

import (
	"github.com/bearbin/go-age"
	"time"
)

// Person type
type Person struct {
	Name             string     `json:"name"`
	FirstName        string     `json:"firstName"`
	LastName         string     `json:"lastName"`
	BirthDate        time.Time  `json:"birthDate"`
	MarriageDate     *time.Time `json:"marriageDate,omitempty"`
	PlaceOfBirth     string     `json:"placeOfBirth"`
	FathersName      string     `json:"fathersName"`
	MothersName      string     `json:"mothersName"`
	MaritalStatus    string     `json:"maritalStatus"`
	SpousesName      string     `json:"spousesName,omitempty"`
	BrothersQuantity int        `json:"brothersQuantity"`
	ChildrenQuantity int        `json:"childrenQuantity"`
	Profession       string     `json:"profession,omitempty"`
	Gender           string     `json:"gender"`
	Contact          *Contact   `json:"contact"`
	Address          *Address   `json:"address"`
}

// GetFullName of a person
func (person Person) GetFullName() string {
	return person.FirstName + " " + person.LastName
}

// Age of a person
func (person Person) Age() int {
	return age.Age(person.BirthDate)
}
