package graphql

import (
	service2 "github.com/BrunoDM2943/church-members-api/internal/service"
	"github.com/graphql-go/graphql"
)

type memberResolver struct {
	service service2.IMemberService
}

func newMemberResolver(service service2.IMemberService) memberResolver {
	return memberResolver{
		service,
	}
}

func (this memberResolver) memberResolver(params graphql.ResolveParams) (interface{}, error) {
	return this.service.FindMembers(params.Args)
}
