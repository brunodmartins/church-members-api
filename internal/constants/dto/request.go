package dto

import (
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/role"
	"time"
)

// RetireMemberRequest for HTTP calls to put member status
// swagger:model RetireMemberRequest
type RetireMemberRequest struct {
	Reason     string    `json:"reason" validate:"required"`
	RetireDate time.Time `json:"date"`
}

// CreateMemberRequest for HTTP calls to post member
// swagger:model CreateMemberRequest
type CreateMemberRequest struct {
	OldChurch              string                `json:"oldChurch"`
	AttendsFridayWorship   *bool                 `json:"attendsFridayWorship" validate:"required"`
	AttendsSaturdayWorship *bool                 `json:"attendsSaturdayWorship" validate:"required"`
	AttendsSundayWorship   *bool                 `json:"attendsSundayWorship" validate:"required"`
	AttendsSundaySchool    *bool                 `json:"attendsSundaySchool" validate:"required"`
	AttendsObservation     string                `json:"attendsObservation"`
	Person                 CreatePersonRequest   `json:"person" validate:"required"`
	Religion               CreateReligionRequest `json:"religion" validate:"required"`
}

func (dto CreateMemberRequest) ToMember() *domain.Member {
	return &domain.Member{
		OldChurch:              dto.OldChurch,
		AttendsFridayWorship:   *dto.AttendsFridayWorship,
		AttendsSaturdayWorship: *dto.AttendsSaturdayWorship,
		AttendsSundayWorship:   *dto.AttendsSundayWorship,
		AttendsSundaySchool:    *dto.AttendsSundaySchool,
		AttendsObservation:     dto.AttendsObservation,
		Person:                 dto.Person.ToPerson(),
		Religion:               dto.Religion.ToReligion(),
	}
}

// CreatePersonRequest for HTTP calls to post a person
// swagger:model CreatePersonRequest
type CreatePersonRequest struct {
	FirstName        string         `json:"firstName" validate:"required"`
	LastName         string         `json:"lastName" validate:"required"`
	BirthDate        Date           `json:"birthDate" validate:"required"`
	MarriageDate     *Date          `json:"marriageDate"`
	PlaceOfBirth     string         `json:"placeOfBirth"`
	FathersName      string         `json:"fathersName"`
	MothersName      string         `json:"mothersName"`
	SpousesName      string         `json:"spousesName"`
	MaritalStatus    string         `json:"maritalStatus" validate:"eq=SINGLE|eq=WIDOW|eq=MARRIED|eq=DIVORCED"`
	BrothersQuantity int            `json:"brothersQuantity"`
	ChildrenQuantity int            `json:"childrenQuantity"`
	Profession       string         `json:"profession"`
	Gender           string         `json:"gender" validate:"required,eq=M|eq=F"`
	Contact          ContactRequest `json:"contact" validate:"required"`
	Address          AddressRequest `json:"address" validate:"required"`
}

func (dto CreatePersonRequest) ToPerson() *domain.Person {
	return &domain.Person{
		Name:             fmt.Sprintf("%s %s", dto.FirstName, dto.LastName),
		FirstName:        dto.FirstName,
		LastName:         dto.LastName,
		BirthDate:        dto.BirthDate.Time,
		MarriageDate:     ToTime(dto.MarriageDate),
		PlaceOfBirth:     dto.PlaceOfBirth,
		FathersName:      dto.FathersName,
		MothersName:      dto.MothersName,
		SpousesName:      dto.SpousesName,
		MaritalStatus:    dto.MaritalStatus,
		BrothersQuantity: dto.BrothersQuantity,
		ChildrenQuantity: dto.ChildrenQuantity,
		Profession:       dto.Profession,
		Gender:           dto.Gender,
		Contact:          dto.Contact.ToContact(),
		Address:          dto.Address.ToAddress(),
	}
}

// ContactRequest for HTTP calls to post a person
// swagger:model ContactRequest
type ContactRequest struct {
	PhoneArea     int    `json:"phoneArea"`
	Phone         int    `json:"phone"`
	CellPhoneArea int    `json:"cellPhoneArea"`
	CellPhone     int    `json:"cellPhone"`
	Email         string `json:"email"`
}

func (dto ContactRequest) ToContact() *domain.Contact {
	return &domain.Contact{
		PhoneArea:     dto.PhoneArea,
		Phone:         dto.Phone,
		CellPhoneArea: dto.CellPhoneArea,
		CellPhone:     dto.CellPhone,
		Email:         dto.Email,
	}
}

// AddressRequest for HTTP calls to post a person
// swagger:model AddressRequest
type AddressRequest struct {
	ZipCode  string `json:"zipCode" validate:"required"`
	State    string `json:"state" validate:"required"`
	City     string `json:"city" validate:"required"`
	Address  string `json:"address" validate:"required"`
	District string `json:"district" validate:"required"`
	Number   int    `json:"number" validate:"required"`
	MoreInfo string `json:"moreInfo"`
}

func (dto AddressRequest) ToAddress() *domain.Address {
	return &domain.Address{
		ZipCode:  dto.ZipCode,
		State:    dto.State,
		City:     dto.City,
		Address:  dto.Address,
		District: dto.District,
		Number:   dto.Number,
		MoreInfo: dto.MoreInfo,
	}
}

// CreateReligionRequest for HTTP calls to post a person
// swagger:model CreateReligionRequest
type CreateReligionRequest struct {
	FathersReligion   string `json:"fathersReligion"`
	BaptismPlace      string `json:"baptismPlace"`
	LearnedGospelAge  int    `json:"learnedGospelAge"`
	AcceptedJesus     bool   `json:"acceptedJesus"`
	Baptized          bool   `json:"baptized"`
	CatholicBaptized  bool   `json:"catholicBaptized"`
	KnowsTithe        bool   `json:"knowsTithe"`
	AgreesTithe       bool   `json:"agreesTithe"`
	Tithe             bool   `json:"tithe"`
	AcceptedJesusDate *Date  `json:"acceptedJesusDate"`
	BaptismDate       *Date  `json:"baptismDate"`
}

func (dto CreateReligionRequest) ToReligion() *domain.Religion {
	return &domain.Religion{
		FathersReligion:   dto.FathersReligion,
		BaptismPlace:      dto.BaptismPlace,
		LearnedGospelAge:  dto.LearnedGospelAge,
		AcceptedJesus:     dto.AcceptedJesus,
		Baptized:          dto.Baptized,
		CatholicBaptized:  dto.CatholicBaptized,
		KnowsTithe:        dto.KnowsTithe,
		AgreesTithe:       dto.AgreesTithe,
		Tithe:             dto.Tithe,
		AcceptedJesusDate: ToTime(dto.AcceptedJesusDate),
		BaptismDate:       ToTime(dto.BaptismDate),
	}
}

// CreateUserRequest for HTTP calls to post user
// swagger:model CreateUserRequest
type CreateUserRequest struct {
	UserName                       string `json:"username" validate:"required,min=3,max=32"`
	Email                          string `json:"email" validate:"required,email,min=3,max=32"`
	Role                           string `json:"role" validate:"required,eq=ADMIN|eq=USER"`
	Password                       string `json:"password" validate:"required,password"`
	Phone                          string `json:"phone" validate:"required"`
	domain.NotificationPreferences `json:"preferences"`
}

func (r CreateUserRequest) ToUser() *domain.User {
	return domain.NewUser(r.UserName, r.Email, r.Password, r.Phone, role.From(r.Role), r.NotificationPreferences)
}

// UpdatePersonRequest for HTTP calls to put a person
// swagger:model UpdatePersonRequest
type UpdatePersonRequest struct {
	FirstName        string `json:"firstName" validate:"required"`
	LastName         string `json:"lastName" validate:"required"`
	BirthDate        Date   `json:"birthDate" validate:"required"`
	MarriageDate     *Date  `json:"marriageDate"`
	SpousesName      string `json:"spousesName"`
	MaritalStatus    string `json:"maritalStatus" validate:"eq=SINGLE|eq=WIDOW|eq=MARRIED|eq=DIVORCED"`
	ChildrenQuantity int    `json:"childrenQuantity"`
}

func (request UpdatePersonRequest) ToPerson() domain.Person {
	return domain.Person{
		FirstName:        request.FirstName,
		LastName:         request.LastName,
		BirthDate:        request.BirthDate.Time,
		MarriageDate:     ToTime(request.MarriageDate),
		SpousesName:      request.SpousesName,
		MaritalStatus:    request.MaritalStatus,
		ChildrenQuantity: request.ChildrenQuantity,
	}
}
