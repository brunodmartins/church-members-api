package graphql

import (
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/graphql-go/graphql"
)

// CreateSchema builds a GraphQL schema
func CreateSchema(service member.Service) graphql.Schema {
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
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						queryFilters := member.QueryBuilder{}
						for key, value := range params.Args {
							queryFilters.AddFilter(key, value)
						}
						return service.SearchMembers(params.Context, queryFilters.ToSpecification())
					},
				},
			},
		}),
	})
	return schema

}
