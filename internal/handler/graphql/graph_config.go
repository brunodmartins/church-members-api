package graphql

import (
	member2 "github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/graphql-go/graphql"
)

//CreateSchema builds a GraphQL schema
func CreateSchema(service member2.Service) graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"member": &graphql.Field{
					Type: graphql.NewList(memberType),
					Args: graphql.FieldConfigArgument{
						"gender": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"active": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: newMemberResolver(service).memberResolver,
				},
			},
		}),
	})
	return schema

}
