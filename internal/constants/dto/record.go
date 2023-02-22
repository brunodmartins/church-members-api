package dto

import (
	"github.com/brunodmartins/church-members-api/internal/constants/enum/role"
	"github.com/brunodmartins/church-members-api/platform/utils"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

// MemberItem for dynamoDB struct
type MemberItem struct {
	ID                     string     `dynamodbav:"id,omitempty"`
	ChurchID               string     `dynamodbav:"church_id"`
	OldChurch              string     `dynamodbav:"oldChurch,omitempty"`
	AttendsFridayWorship   bool       `dynamodbav:"attendsFridayWorship"`
	AttendsSaturdayWorship bool       `dynamodbav:"attendsSaturdayWorship"`
	AttendsSundayWorship   bool       `dynamodbav:"attendsSundayWorship"`
	AttendsSundaySchool    bool       `dynamodbav:"attendsSundaySchool"`
	AttendsObservation     string     `dynamodbav:"attendsObservation,omitempty"`
	Name                   string     `dynamodbav:"name"`
	FirstName              string     `dynamodbav:"firstName"`
	LastName               string     `dynamodbav:"lastName"`
	BirthDate              time.Time  `dynamodbav:"birthDate"`
	MarriageDate           *time.Time `dynamodbav:"marriageDate,omitempty"`
	PlaceOfBirth           string     `dynamodbav:"placeOfBirth"`
	FathersName            string     `dynamodbav:"fathersName"`
	MothersName            string     `dynamodbav:"mothersName"`
	SpousesName            string     `dynamodbav:"spousesName,omitempty"`
	BrothersQuantity       int        `dynamodbav:"brothersQuantity"`
	ChildrenQuantity       int        `dynamodbav:"childrensQuantity"`
	Profession             string     `dynamodbav:"profession,omitempty"`
	Gender                 string     `dynamodbav:"gender"`
	PhoneArea              int        `dynamodbav:"phoneArea,omitempty"`
	Phone                  int        `dynamodbav:"phone,omitempty"`
	CellPhoneArea          int        `dynamodbav:"cellPhoneArea"`
	CellPhone              int        `dynamodbav:"cellPhone"`
	Email                  string     `dynamodbav:"email"`
	ZipCode                string     `dynamodbav:"zipCode"`
	State                  string     `dynamodbav:"state"`
	City                   string     `dynamodbav:"city"`
	Address                string     `dynamodbav:"address"`
	District               string     `dynamodbav:"district"`
	AddressNumber          int        `dynamodbav:"addressNumber"`
	MoreInfo               string     `dynamodbav:"moreInfo"`
	FathersReligion        string     `dynamodbav:"fathersReligion,omitempty"`
	BaptismPlace           string     `dynamodbav:"baptismPlace"`
	LearnedGospelAge       int        `dynamodbav:"learnedGospelAge"`
	AcceptedJesus          bool       `dynamodbav:"acceptedJesus"`
	Baptized               bool       `dynamodbav:"baptized"`
	CatholicBaptized       bool       `dynamodbav:"catholicBaptized"`
	KnowsTithe             bool       `dynamodbav:"knowsTithe"`
	AgreesTithe            bool       `dynamodbav:"agreesTithe"`
	Tithe                  bool       `dynamodbav:"tithe"`
	AcceptedJesusDate      *time.Time `dynamodbav:"acceptedJesusDate"`
	BaptismDate            *time.Time `dynamodbav:"baptismDate"`
	Active                 bool       `dynamodbav:"active,omitempty"`
	BirthDateShort         string     `dynamodbav:"birthDateShort"`
	MarriageDateShort      string     `dynamodbav:"marriageDateShort"`
}

// NewMemberItem creates a MemberItem from a domain.Member
func NewMemberItem(member *domain.Member) *MemberItem {
	return &MemberItem{
		ID:                     member.ID,
		ChurchID:               member.ChurchID,
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
		ChildrenQuantity:       member.Person.ChildrenQuantity,
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
		AddressNumber:          member.Person.Address.Number,
		MoreInfo:               member.Person.Address.MoreInfo,
		FathersReligion:        member.Religion.FathersReligion,
		BaptismPlace:           member.Religion.BaptismPlace,
		LearnedGospelAge:       member.Religion.LearnedGospelAge,
		AcceptedJesus:          member.Religion.AcceptedJesus,
		Baptized:               member.Religion.Baptized,
		CatholicBaptized:       member.Religion.CatholicBaptized,
		KnowsTithe:             member.Religion.KnowsTithe,
		AgreesTithe:            member.Religion.AgreesTithe,
		Tithe:                  member.Religion.Tithe,
		AcceptedJesusDate:      member.Religion.AcceptedJesusDate,
		BaptismDate:            member.Religion.BaptismDate,
		Active:                 member.Active,
		BirthDateShort:         utils.ConvertDate(member.Person.BirthDate),
		MarriageDateShort:      convertMarriageDate(member.Person.MarriageDate),
	}
}

func convertMarriageDate(marriageDate *time.Time) string {
	if marriageDate != nil {
		return utils.ConvertDate(*marriageDate)
	}
	return ""
}

// ToMember converts a MemberItem into a domain.Member
func (item *MemberItem) ToMember() *domain.Member {
	return &domain.Member{
		ID:                     item.ID,
		ChurchID:               item.ChurchID,
		OldChurch:              item.OldChurch,
		AttendsFridayWorship:   item.AttendsFridayWorship,
		AttendsSaturdayWorship: item.AttendsSaturdayWorship,
		AttendsSundayWorship:   item.AttendsSundayWorship,
		AttendsSundaySchool:    item.AttendsSundaySchool,
		AttendsObservation:     item.AttendsObservation,
		Person: domain.Person{
			Name:             item.Name,
			FirstName:        item.FirstName,
			LastName:         item.LastName,
			BirthDate:        item.BirthDate,
			MarriageDate:     item.MarriageDate,
			PlaceOfBirth:     item.PlaceOfBirth,
			FathersName:      item.FathersName,
			MothersName:      item.MothersName,
			SpousesName:      item.SpousesName,
			BrothersQuantity: item.BrothersQuantity,
			ChildrenQuantity: item.ChildrenQuantity,
			Profession:       item.Profession,
			Gender:           item.Gender,
			Contact: domain.Contact{
				PhoneArea:     item.PhoneArea,
				Phone:         item.Phone,
				CellPhoneArea: item.CellPhoneArea,
				CellPhone:     item.CellPhone,
				Email:         item.Email,
			},
			Address: domain.Address{
				ZipCode:  item.ZipCode,
				State:    item.State,
				City:     item.City,
				Address:  item.Address,
				District: item.District,
				Number:   item.AddressNumber,
				MoreInfo: item.MoreInfo,
			},
		},
		Religion: domain.Religion{
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

type UserItem struct {
	ID              string `dynamodbav:"id"`
	ChurchID        string `dynamodbav:"church_id"`
	UserName        string `dynamodbav:"username"`
	Email           string `dynamodbav:"email"`
	Role            string `dynamodbav:"role"`
	Password        string `dynamodbav:"password"`
	Phone           string `dynamodbav:"phone"`
	SendDailySMS    bool   `dynamodbav:"send_daily_sms"`
	SendWeeklyEmail bool   `dynamodbav:"send_weekly_email"`
}

// NewUserItem creates a UserItem from a domain.User
func NewUserItem(user *domain.User) *UserItem {
	return &UserItem{
		ID:              user.ID,
		UserName:        user.UserName,
		Email:           user.Email,
		Role:            user.Role.String(),
		Password:        string(user.Password),
		Phone:           user.Phone,
		SendDailySMS:    user.Preferences.SendDailySMS,
		SendWeeklyEmail: user.Preferences.SendWeeklyEmail,
	}
}

// ToUser converts a UserItem into a domain.User
func (item *UserItem) ToUser() *domain.User {
	return &domain.User{
		ID:       item.ID,
		ChurchID: item.ChurchID,
		UserName: item.UserName,
		Email:    item.Email,
		Role:     role.From(item.Role),
		Password: []byte(item.Password),
		Phone:    item.Phone,
		Preferences: domain.NotificationPreferences{
			SendDailySMS:    item.SendDailySMS,
			SendWeeklyEmail: item.SendWeeklyEmail,
		},
	}
}
