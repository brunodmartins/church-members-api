package domain

import (
	"github.com/bearbin/go-age"
	"github.com/brunodmartins/church-members-api/internal/constants/enum"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/classification"
	"time"
)

// Member struct
type Member struct {
	ID                     string     `json:"id"`
	ChurchID               string     `json:"church_id"`
	OldChurch              string     `json:"oldChurch,omitempty"`
	AttendsFridayWorship   bool       `json:"attendsFridayWorship"`
	AttendsSaturdayWorship bool       `json:"attendsSaturdayWorship"`
	AttendsSundayWorship   bool       `json:"attendsSundayWorship"`
	AttendsSundaySchool    bool       `json:"attendsSundaySchool"`
	AttendsObservation     string     `json:"attendsObservation,omitempty"`
	Person                 *Person    `json:"person"`
	Religion               *Religion  `json:"religion"`
	Active                 bool       `json:"active,omitempty"`
	MembershipStartDate    time.Time  `json:"membership_start_date"`
	MembershipEndDate      *time.Time `json:"membership_end_date,omitempty"`
	MembershipEndReason    string     `json:"membership_end_reason,omitempty"`
}

// Classification returns a member classification based on age and marriage
func (member Member) Classification() enum.Classification {
	personAge := age.Age(member.Person.BirthDate)
	if personAge < 15 {
		return classification.CHILDREN
	} else if personAge < 18 {
		return classification.TEEN
	} else if personAge < 30 && member.Person.MarriageDate == nil {
		return classification.YOUNG
	} else {
		return classification.ADULT
	}
}

// IsLegal validate if a member is legal only if it's not a children
func (member Member) IsLegal() bool {
	return member.Classification() != classification.CHILDREN
}

func (member Member) MembershipEndCurrentYear() bool {
	return member.MembershipEndDate != nil && member.MembershipEndDate.Year() == time.Now().Year()
}
