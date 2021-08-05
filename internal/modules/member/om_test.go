package member_test

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
)

var (
	genericError = errors.New("generic error")
)

func BuildMembers(size int) []*domain.Member {
	var members []*domain.Member
	for i:=0;i<size;i++ {
		members = append(members, buildMember(domain.NewID()))
	}
	return members
}

func buildMember(id string) *domain.Member{
	return &domain.Member{
		ID: id,
	}
}
