package graphql

import (
	"github.com/BrunoDM2943/church-members-api/member"
	"github.com/graphql-go/graphql"
)

func CreateSchema(service member.Service) graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"member": &graphql.Field{
					Type: graphql.NewList(memberType),
					Args: graphql.FieldConfigArgument{
						"sexo": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"active": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
					},
					Resolve: newMemberResolver(service).memberResolver,
				},
			},
		}),
	})
	return schema

}
