package dto

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"time"
)

//MemberItem for dynamoDB struct
type MemberItem struct {
	ID                     string    `dynamodbav:"id,omitempty"`
	OldChurch              string    `dynamodbav:"oldChurch,omitempty"`
	AttendsFridayWorship   bool      `dynamodbav:"attendsFridayWorship"`
	AttendsSaturdayWorship bool      `dynamodbav:"attendsSaturdayWorship"`
	AttendsSundayWorship   bool      `dynamodbav:"attendsSundayWorship"`
	AttendsSundaySchool    bool      `dynamodbav:"attendsSundaySchool"`
	AttendsObservation     string    `dynamodbav:"attendsObservation,omitempty"`
	Name                   string    `dynamodbav:"name"`
	FirstName              string    `dynamodbav:"firstName"`
	LastName               string    `dynamodbav:"lastName"`
	BirthDate              *time.Time `dynamodbav:"birthDate"`
	MarriageDate           *time.Time `dynamodbav:"marriageDate,omitempty"`
	PlaceOfBirth           string    `dynamodbav:"placeOfBirth"`
	FathersName            string    `dynamodbav:"fathersName"`
	MothersName            string    `dynamodbav:"mothersName"`
	SpousesName            string    `dynamodbav:"spousesName,omitempty"`
	BrothersQuantity       int       `dynamodbav:"brothersQuantity"`
	ChildrensQuantity      int       `dynamodbav:"childrensQuantity"`
	Profession             string    `dynamodbav:"profession,omitempty"`
	Gender                 string    `dynamodbav:"gender"`
	PhoneArea              int       `dynamodbav:"phoneArea,omitempty"`
	Phone                  int       `dynamodbav:"phone,omitempty"`
	CellPhoneArea          int       `dynamodbav:"cellPhoneArea"`
	CellPhone              int       `dynamodbav:"cellPhone"`
	Email                  string    `dynamodbav:"email"`
	ZipCode                string    `dynamodbav:"zipCode"`
	State                  string    `dynamodbav:"state"`
	City                   string    `dynamodbav:"city"`
	Address                string    `dynamodbav:"address"`
	District               string    `dynamodbav:"district"`
	AddressNumber          int       `dynamodbav:"addressNumber"`
	MoreInfo               string    `dynamodbav:"moreInfo"`
	FathersReligion        string    `dynamodbav:"fathersReligion,omitempty"`
	BaptismPlace           string    `dynamodbav:"baptismPlace"`
	LearnedGospelAge       int       `dynamodbav:"learnedGospelAge"`
	AcceptedJesus          bool      `dynamodbav:"acceptedJesus"`
	Baptized               bool      `dynamodbav:"baptized"`
	CatholicBaptized       bool      `dynamodbav:"catholicBaptized"`
	KnowsTithe             bool      `dynamodbav:"knowsTithe"`
	AgreesTithe            bool      `dynamodbav:"agreesTithe"`
	Tithe                  bool      `dynamodbav:"tithe"`
	AcceptedJesusDate      *time.Time `dynamodbav:"acceptedJesusDate"`
	BaptismDate            *time.Time `dynamodbav:"baptismDate"`
	Active                 bool      `dynamodbav:"active,omitempty"`
}

func NewMemberItem(member *model.Member) *MemberItem {
	return &MemberItem{
		ID:                     member.ID,
		OldChurch:              member.OldChurch,
		AttendsFridayWorship:   member.AttendsFridayWorship,
		AttendsSaturdayWorship: member.AttendsSaturdayWorship,
		AttendsSundayWorship:   member.AttendsSundayWorship,
		AttendsSundaySchool:    member.AttendsSundaySchool,
		AttendsObservation:     member.AttendsObservation,
		Name:                   member.Person.GetFullName(),
		FirstName:              member.Person.FirstName,
		LastName:               member.Person.LastName,
		BirthDate:              member.Person.BirthDate,
		MarriageDate:           member.Person.MarriageDate,
		PlaceOfBirth:           member.Person.PlaceOfBirth,
		FathersName:            member.Person.FathersName,
		MothersName:            member.Person.MothersName,
		SpousesName:            member.Person.SpousesName,
		BrothersQuantity:       member.Person.BrothersQuantity,
		ChildrensQuantity:      member.Person.ChildrensQuantity,
		Profession:             member.Person.Profession,
		Gender:                 member.Person.Gender,
		PhoneArea:              member.Person.Contact.PhoneArea,
		Phone:                  member.Person.Contact.Phone,
		CellPhoneArea:          member.Person.Contact.CellPhoneArea,
		CellPhone:              member.Person.Contact.CellPhone,
		Email:                  member.Person.Contact.Email,
		ZipCode:                member.Person.Address.ZipCode,
		State:                  member.Person.Address.State,
		City:                   member.Person.Address.City,
		Address:                member.Person.Address.Address,
		District:               member.Person.Address.District,
		AddressNumber:                 member.Person.Address.Number,
		MoreInfo:               member.Person.Address.MoreInfo,

		FathersReligion:   member.Religion.FathersReligion,
		BaptismPlace:      member.Religion.BaptismPlace,
		LearnedGospelAge:  member.Religion.LearnedGospelAge,
		AcceptedJesus:     member.Religion.AcceptedJesus,
		Baptized:          member.Religion.Baptized,
		CatholicBaptized:  member.Religion.CatholicBaptized,
		KnowsTithe:        member.Religion.KnowsTithe,
		AgreesTithe:       member.Religion.AgreesTithe,
		Tithe:             member.Religion.Tithe,
		AcceptedJesusDate: member.Religion.AcceptedJesusDate,
		BaptismDate:       member.Religion.BaptismDate,
		Active:            member.Active,
	}
}

func (item *MemberItem) ToMember() *model.Member {
	return &model.Member{
		ID:                     item.ID,
		OldChurch:              item.OldChurch,
		AttendsFridayWorship:   item.AttendsFridayWorship,
		AttendsSaturdayWorship: item.AttendsSaturdayWorship,
		AttendsSundayWorship:   item.AttendsSundayWorship,
		AttendsSundaySchool:    item.AttendsSundaySchool,
		AttendsObservation:     item.AttendsObservation,
		Person: model.Person{
			Name:              item.Name,
			FirstName:         item.FirstName,
			LastName:          item.LastName,
			BirthDate:         item.BirthDate,
			MarriageDate:      item.MarriageDate,
			PlaceOfBirth:      item.PlaceOfBirth,
			FathersName:       item.FathersName,
			MothersName:       item.MothersName,
			SpousesName:       item.SpousesName,
			BrothersQuantity:  item.BrothersQuantity,
			ChildrensQuantity: item.ChildrensQuantity,
			Profession:        item.Profession,
			Gender:            item.Gender,
			Contact: model.Contact{
				PhoneArea:     item.PhoneArea,
				Phone:         item.Phone,
				CellPhoneArea: item.CellPhoneArea,
				CellPhone:     item.CellPhone,
				Email:         item.Email,
			},
			Address: model.Address{
				ZipCode:  item.ZipCode,
				State:    item.State,
				City:     item.City,
				Address:  item.Address,
				District: item.District,
				Number:   item.AddressNumber,
				MoreInfo: item.MoreInfo,
			},
		},
		Religion: model.Religion{
			FathersReligion:   item.FathersReligion,
			BaptismPlace:      item.BaptismPlace,
			LearnedGospelAge:  item.LearnedGospelAge,
			AcceptedJesus:     item.AcceptedJesus,
			Baptized:          item.Baptized,
			CatholicBaptized:  item.CatholicBaptized,
			KnowsTithe:        item.KnowsTithe,
			AgreesTithe:       item.AgreesTithe,
			Tithe:             item.Tithe,
			AcceptedJesusDate: item.AcceptedJesusDate,
			BaptismDate:       item.BaptismDate,
		},
		Active: item.Active,
	}
}
