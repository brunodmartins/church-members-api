package graphql

import (
	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/bearbin/go-age"
	"github.com/graphql-go/graphql"
)

var memberType = graphql.NewObject(graphql.ObjectConfig{
	Name: "member",
	Fields: graphql.Fields{
		"id": &graphql.Field{
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
		"classificacao": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				membro := p.Source.(*entity.Membro)
				return membro.Classificacao(), nil
			},
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
		"dtNascimento": &graphql.Field{
			Type: graphql.DateTime,
		},
		"dtCasamento": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				pessoa := p.Source.(entity.Pessoa)
				if pessoa.DtCasamento.IsZero() {
					return nil, nil
				}
				return pessoa.DtCasamento, nil
			},
		},
		"contato": &graphql.Field{
			Type: contactType,
		},
		"nomeConjuge": &graphql.Field{
			Type: graphql.String,
		},
		"endereco": &graphql.Field{
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
				contato := p.Source.(entity.Contato)
				if contato.Celular != 0 {
					return contato.GetFormattedCellPhone(), nil
				}
				return nil, nil
			},
		},
		"phone": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				contato := p.Source.(entity.Contato)
				if contato.Telefone != 0 {
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
		"cep": &graphql.Field{
			Type: graphql.String,
		},
		"uf": &graphql.Field{
			Type: graphql.String,
		},
		"cidade": &graphql.Field{
			Type: graphql.String,
		},
		"logradouro": &graphql.Field{
			Type: graphql.String,
		},
		"bairro": &graphql.Field{
			Type: graphql.String,
		},
		"numero": &graphql.Field{
			Type: graphql.Int,
		},
		"full": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				endereco := p.Source.(entity.Endereco)
				return endereco.GetFormatted(), nil
			},
		},
	},
})
