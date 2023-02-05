package graphql

import (
	"github.com/graphql-go/graphql"
)

var memberType = graphql.NewObject(graphql.ObjectConfig{
	Name: "member",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:    graphql.String,
			Resolve: idResolver,
		},
		"active": &graphql.Field{
			Type: graphql.Boolean,
		},
		"person": &graphql.Field{
			Type: personType,
		},
		"classification": &graphql.Field{
			Type:    graphql.String,
			Resolve: classificationResolver,
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
			Type:    graphql.String,
			Resolve: fullNameResolver,
		},
		"gender": &graphql.Field{
			Type: graphql.String,
		},
		"age": &graphql.Field{
			Type:    graphql.Int,
			Resolve: ageResolver,
		},
		"birthDate": &graphql.Field{
			Type: graphql.DateTime,
		},
		"marriageDate": &graphql.Field{
			Type:    graphql.DateTime,
			Resolve: marriageDateResolver,
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
			Type:    graphql.String,
			Resolve: cellPhoneResolver,
		},
		"phone": &graphql.Field{
			Type:    graphql.String,
			Resolve: phoneResolver,
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
			Type:    graphql.String,
			Resolve: fullAddressResolver,
		},
	},
})
