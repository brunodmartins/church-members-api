package graphql

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/bearbin/go-age"
	"github.com/graphql-go/graphql"
)

func idResolver(p graphql.ResolveParams) (interface{}, error) {
	member := p.Source.(*domain.Member)
	return member.ID, nil
}

func classificationResolver(p graphql.ResolveParams) (interface{}, error) {
	member := p.Source.(*domain.Member)
	return member.Classification().String(), nil
}

func fullNameResolver(p graphql.ResolveParams) (i interface{}, e error) {
	person := p.Source.(domain.Person)
	return person.GetFullName(), nil
}

func ageResolver(p graphql.ResolveParams) (i interface{}, e error) {
	return age.Age(p.Source.(domain.Person).BirthDate), nil
}

func marriageDateResolver(p graphql.ResolveParams) (interface{}, error) {
	person := p.Source.(domain.Person)
	if person.MarriageDate == nil {
		return nil, nil
	}
	return person.MarriageDate, nil
}

func cellPhoneResolver(p graphql.ResolveParams) (interface{}, error) {
	contact := p.Source.(domain.Contact)
	return contact.GetFormattedCellPhone(), nil
}

func phoneResolver(p graphql.ResolveParams) (interface{}, error) {
	contact := p.Source.(domain.Contact)
	return contact.GetFormattedPhone(), nil
}

func fullAddressResolver(p graphql.ResolveParams) (interface{}, error) {
	address := p.Source.(domain.Address)
	return address.String(), nil
}