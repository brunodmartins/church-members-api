package graphql

import (
	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/bearbin/go-age"
	"github.com/graphql-go/graphql"
)

var memberType = graphql.NewObject(graphql.ObjectConfig{
	Name: "member",
	Fields: graphql.Fields{
		"_id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				member := p.Source.(*entity.Membro)
				return member.ID.String(), nil
			},
		},
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
		"fullName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				pessoa := p.Source.(entity.Pessoa)
				return pessoa.GetFullName(), nil
			},
		},
		"sexo": &graphql.Field{
			Type: graphql.String,
		},
		"age": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				return age.Age(p.Source.(entity.Pessoa).DtNascimento), nil
			},
		},
	},
})
