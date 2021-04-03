package model

import (
	"time"
)

//Person type
type Person struct {
	FirstName         string    `json:"firstName" bson:"firstName"`
	LastName          string    `json:"lastName" bson:"lastName"`
	BirthDate         time.Time `json:"birthDate" bson:"birthDate"`
	MarriageDate      time.Time `json:"marriageDate,omitempty" bson:"marriageDate"`
	PlaceOfBirth      string    `json:"placeOfBirth" bson:"placeOfBirth"`
	FathersName       string    `json:"fathersName" bson:"fathersName"`
	MothersName       string    `json:"mothersName" bson:"mothersName"`
	SpousesName       string    `json:"spousesName,omitempty" bson:"spousesName"`
	BrothersQuantity  int       `json:"brothersQuantity" bson:"brothersQuantity"`
	ChildrensQuantity int       `json:"childrensQuantity" bson:"childrensQuantity"`
	Profession        string    `json:"profession,omitempty"`
	Gender            string    `json:"gender"`
	Contact           Contact   `json:"contact"`
	Address           Address   `json:"address"`
}

//GetFullName of a person
func (person Person) GetFullName() string {
	return person.FirstName + " " + person.LastName
}
