package entity

import (
	"github.com/bearbin/go-age"
)

//Member struct
type Member struct {
	ID                     ID       `json:"id" bson:"_id,omitempty"`
	OldChurch              string   `json:"oldChurch,omitempty"`
	AttendsFridayWorship   bool     `json:"attendsFridayWorship" bson:"attendsFridayWorship"`
	AttendsSaturdayWorship bool     `json:"attendsSaturdayWorship" bson:"attendsSaturdayWorship"`
	AttendsSundayWorship   bool     `json:"attendsSundayWorship" bson:"attendsSundayWorship"`
	AttendsSundaySchool    bool     `json:"attendsSundaySchool" bson:"attendsSundaySchool"`
	AttendsObservation     string   `json:"attendsObservation,omitempty" bson:"attendsObservation"`
	Person                 Person   `json:"person"`
	Religion               Religion `json:"religion"`
	Active                 bool     `json:"active,omitempty" bson:"active"`
}

//Classification returns a member classification based on age and marriage
func (member Member) Classification() string {
	age := age.Age(member.Person.BirthDate)
	if age < 15 {
		return "Criança"
	} else if age < 18 {
		return "Adolescente"
	} else if age < 30 && member.Person.MarriageDate.IsZero() {
		return "Jovem"
	} else {
		return "Adulto"
	}
}
