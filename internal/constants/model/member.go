package model

import (
	"regexp"

	"github.com/bearbin/go-age"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
func (member Member) Classification() string {
	age := age.Age(*member.Person.BirthDate)
	if age < 15 {
		return "Children"
	} else if age < 18 {
		return "Teen"
	} else if age < 30 && member.Person.MarriageDate == nil {
		return "Young"
	} else {
		return "Adult"
	}
}

//Classification returns a member classification based on age and marriage
func (member Member) ClassificationLocalized(localizer *i18n.Localizer) string {
	classification := member.Classification()
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Classification." + classification,
			Other: classification,
		},
	})
}

func IsValidID(id string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(id)
}

func NewID() string {
	return uuid.NewString()
}
