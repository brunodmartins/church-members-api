package graphql

import (
	"github.com/BrunoDM2943/church-members-api/member/service"
	"github.com/graphql-go/graphql"
)

type memberResolver struct {
	service service.IMemberService
}

func newMemberResolver(service service.IMemberService) memberResolver {
	return memberResolver{
		service,
	}
}

func (this memberResolver) memberResolver(params graphql.ResolveParams) (interface{}, error) {
	return this.service.FindMembers(params.Args)
}
