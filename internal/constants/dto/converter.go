package dto

import (
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

// NewGetMemberResponse builds a member response
func NewGetMemberResponse(member *domain.Member) *GetMemberResponse {
	result := new(GetMemberResponse)
	result.ID = member.ID
	result.Active = member.Active
	result.Classification = member.Classification().String()
	result.Person = buildPersonResponse(member.Person)
	result.Baptism = buildBaptismResponse(member.Religion)
	return result
}

func buildPersonResponse(person *domain.Person) *GetPersonResponse {
	return &GetPersonResponse{
		FirstName:        person.FirstName,
		LastName:         person.LastName,
		FullName:         person.GetFullName(),
		Gender:           person.Gender,
		Age:              person.Age(),
		BirthDate:        person.BirthDate,
		MarriageDate:     person.MarriageDate,
		SpousesName:      person.SpousesName,
		MaritalStatus:    person.MaritalStatus,
		ChildrenQuantity: person.ChildrenQuantity,
		Contact:          buildContactResponse(person.Contact),
		Address:          buildAddressResponse(person.Address),
	}
}

func buildContactResponse(contact *domain.Contact) *GetContactResponse {
	if contact == nil {
		return nil
	}
	return &GetContactResponse{
		CellPhoneArea:      contact.CellPhoneArea,
		CellPhone:          contact.CellPhone,
		PhoneArea:          contact.PhoneArea,
		Phone:              contact.Phone,
		CellPhoneFormatted: contact.GetFormattedCellPhone(),
		PhoneFormatted:     contact.GetFormattedPhone(),
		Email:              contact.Email,
	}
}

func buildAddressResponse(address *domain.Address) *GetAddressResponse {
	if address == nil {
		return nil
	}
	return &GetAddressResponse{
		ZipCode:  address.ZipCode,
		State:    address.State,
		City:     address.City,
		Address:  address.Address,
		District: address.District,
		Number:   address.Number,
		MoreInfo: address.MoreInfo,
		Full:     address.String(),
	}
}

func buildBaptismResponse(religion *domain.Religion) *GetBaptismResponse {
	return &GetBaptismResponse{
		BaptismPlace:      religion.BaptismPlace,
		AcceptedJesus:     religion.AcceptedJesus,
		Baptized:          religion.Baptized,
		CatholicBaptized:  religion.CatholicBaptized,
		AcceptedJesusDate: religion.AcceptedJesusDate,
		BaptismDate:       religion.BaptismDate,
	}
}
