package graphql

import(
	"github.com/BrunoDM2943/church-members-api/member"
	"github.com/graphql-go/graphql"
)

func CreateSchema(service member.Service) graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"person": &graphql.Field{
					Type: graphql.NewList(personType),
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return nil, nil
					},
				},
				"member":  &graphql.Field{
					Type: graphql.NewList(memberType),
					Resolve: newMemberResolver(service).memberResolver,
				},
			},
		}),
	})
	return schema

}
