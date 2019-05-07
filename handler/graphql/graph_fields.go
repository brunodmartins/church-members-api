package graphql

import "github.com/graphql-go/graphql"

var memberType = graphql.NewObject(graphql.ObjectConfig{
	Name: "member",
	Fields: graphql.Fields{
		"active": &graphql.Field{
			Type: graphql.Boolean,
		},
		"pessoa": &graphql.Field{
			Type: personType,
		},
	},
})

var personType = graphql.NewObject(graphql.ObjectConfig{
	Name: "person",
	Fields: graphql.Fields{
		"nome": &graphql.Field{
			Type: graphql.String,
		},
		"sobrenome": &graphql.Field{
			Type: graphql.String,
		},
		"sexo": &graphql.Field{
			Type: graphql.String,
		},
	},
})
