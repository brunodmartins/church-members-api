package model

import (
	"github.com/BrunoDM2943/church-members-api/internal/infra/i18n"
	"testing"
	"time"
	"gopkg.in/mgo.v2/bson"
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

func TestClassification(t *testing.T) {
	i18n.BootStrapI18N()
	t.Run("Children", func(t *testing.T) {
		assert.Equal(t, "Children", Member{
			Person: Person{
				BirthDate: time.Now(),
			},
		}.Classification())
	})
	t.Run("Teen", func(t *testing.T) {
		assert.Equal(t, "Teen", Member{
			Person: Person{
				BirthDate: time.Now().AddDate(-17, 0, 0),
			},
		}.Classification())
	})
	t.Run("Young", func(t *testing.T) {
		assert.Equal(t, "Young", Member{
			Person: Person{
				BirthDate: time.Now().AddDate(-29, 0, 0),
			},
		}.Classification())
	})
	t.Run("Adult Single", func(t *testing.T) {
		assert.Equal(t, "Adult", Member{
			Person: Person{
				BirthDate: time.Now().AddDate(-33, 0, 0),
			},
		}.Classification())
	})
	t.Run("Adult Married", func(t *testing.T) {
		assert.Equal(t, "Adult", Member{
			Person: Person{
				BirthDate:    time.Now().AddDate(-25, 0, 0),
				MarriageDate: time.Now(),
			},
		}.Classification())
	})
	t.Run("Localized", func(t *testing.T) {

		assert.Equal(t, "Adult", Member{
			Person: Person{
				BirthDate:    time.Now().AddDate(-25, 0, 0),
				MarriageDate: time.Now(),
			},
		}.ClassificationLocalized(i18n.Localizer))
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

func TestGetFullName(t *testing.T) {
	assert.Equal(t, Person{
		FirstName: "John",
		LastName:  "Doe",
	}.GetFullName(), "John Doe")
}

func TestID(t *testing.T) {
	t.Run("Test Marshal and Unmarshal", func(t *testing.T) {
		originalId := NewID()
		bytes, _ := originalId.MarshalJSON()
		newId := NewID()
		newId.UnmarshalJSON(bytes)
		assert.Equal(t, originalId.String(), newId.String())
	})
	t.Run("Test Get BSON", func(t *testing.T) {
		originalId := NewID()
		bsonValue, _ := originalId.GetBSON()
		objectId := bsonValue.(bson.ObjectId)
		assert.Equal(t, originalId.String(), objectId.Hex())
	})

}