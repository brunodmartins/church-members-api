package graphql

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/bearbin/go-age"
	"github.com/graphql-go/graphql"
)

var memberType = graphql.NewObject(graphql.ObjectConfig{
	Name: "member",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				member := p.Source.(*model.Member)
				return member.ID.String(), nil
			},
		},
		"active": &graphql.Field{
			Type: graphql.Boolean,
		},
		"person": &graphql.Field{
			Type: personType,
		},
		"classification": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				member := p.Source.(*model.Member)
				return member.Classification(), nil
			},
		},
	},
})

var personType = graphql.NewObject(graphql.ObjectConfig{
	Name: "person",
	Fields: graphql.Fields{
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
		"fullName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				person := p.Source.(model.Person)
				return person.GetFullName(), nil
			},
		},
		"gender": &graphql.Field{
			Type: graphql.String,
		},
		"age": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				return age.Age(p.Source.(model.Person).BirthDate), nil
			},
		},
		"birthDate": &graphql.Field{
			Type: graphql.DateTime,
		},
		"marriageDate": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				person := p.Source.(model.Person)
				if person.MarriageDate.IsZero() {
					return nil, nil
				}
				return person.MarriageDate, nil
			},
		},
		"contact": &graphql.Field{
			Type: contactType,
		},
		"spousesName": &graphql.Field{
			Type: graphql.String,
		},
		"address": &graphql.Field{
			Type: addressType,
		},
	},
})

var contactType = graphql.NewObject(graphql.ObjectConfig{
	Name: "contact",
	Fields: graphql.Fields{
		"cellphone": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				contato := p.Source.(model.Contact)
				if contato.CellPhone != 0 {
					return contato.GetFormattedCellPhone(), nil
				}
				return nil, nil
			},
		},
		"phone": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				contato := p.Source.(model.Contact)
				if contato.Phone != 0 {
					return contato.GetFormattedPhone(), nil
				}
				return nil, nil
			},
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var addressType = graphql.NewObject(graphql.ObjectConfig{
	Name: "address",
	Fields: graphql.Fields{
		"zipcode": &graphql.Field{
			Type: graphql.String,
		},
		"state": &graphql.Field{
			Type: graphql.String,
		},
		"city": &graphql.Field{
			Type: graphql.String,
		},
		"address": &graphql.Field{
			Type: graphql.String,
		},
		"district": &graphql.Field{
			Type: graphql.String,
		},
		"number": &graphql.Field{
			Type: graphql.Int,
		},
		"full": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				address := p.Source.(model.Address)
				return address.GetFormatted(), nil
			},
		},
	},
})
