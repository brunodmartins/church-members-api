package graphql

import (
	"github.com/BrunoDM2943/church-members-api/member"
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

func (this memberResolver) memberResolver(params graphql.ResolveParams) (interface{}, error) {
	return this.service.FindAll()
}