package domain

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"github.com/bearbin/go-age"
)

//Member struct
type Member struct {
	ID                     string   `json:"id"`
	OldChurch              string   `json:"oldChurch,omitempty"`
	AttendsFridayWorship   bool     `json:"attendsFridayWorship"`
	AttendsSaturdayWorship bool     `json:"attendsSaturdayWorship"`
	AttendsSundayWorship   bool     `json:"attendsSundayWorship"`
	AttendsSundaySchool    bool     `json:"attendsSundaySchool"`
	AttendsObservation     string   `json:"attendsObservation,omitempty"`
	Person                 Person   `json:"person"`
	Religion               Religion `json:"religion"`
	Active                 bool     `json:"active,omitempty"`
}

//Classification returns a member classification based on age and marriage
func (member Member) Classification() enum.Classification {
	age := age.Age(*member.Person.BirthDate)
	if age < 15 {
		return enum.CHILDREN
	} else if age < 18 {
		return enum.TEEN
	} else if age < 30 && member.Person.MarriageDate == nil {
		return enum.YOUNG
	} else {
		return enum.ADULT
	}
}