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
	for key, value := range params.Args {
		queryFilters.AddFilter(key, value)
	}
	return resolver.service.SearchMembers(queryFilters.ToSpecification())
}
