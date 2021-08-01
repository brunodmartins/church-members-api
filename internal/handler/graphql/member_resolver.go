package graphql

import (
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/graphql-go/graphql"
)

type memberResolver struct {
	service member.Service
}

func newMemberResolver(service member.Service) memberResolver {
	return memberResolver{
		service,
	}
}

func (resolver memberResolver) memberResolver(params graphql.ResolveParams) (interface{}, error) {
	queryFilters := member.QuerySpecification{}

	if sex := params.Args["gender"]; sex != nil {
		queryFilters.AddFilter("person.gender", sex)
	}

	if active := params.Args["active"]; active != nil {
		queryFilters.AddFilter("active", active.(bool))
	}

	if name := params.Args["name"]; name != nil {
		queryFilters.AddFilter("name", name)
	}
	return resolver.service.SearchMembers(queryFilters.ToSpecification())
}
