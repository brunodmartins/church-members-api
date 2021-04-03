package model

import (
	"github.com/bearbin/go-age"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//Member struct
type Member struct {
	ID                     ID       `json:"id" bson:"_id,omitempty"`
	OldChurch              string   `json:"oldChurch,omitempty" bson:"oldChurch"`
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
		return "Children"
	} else if age < 18 {
		return "Teen"
	} else if age < 30 && member.Person.MarriageDate.IsZero() {
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
