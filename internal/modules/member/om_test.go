package member_test

import (
	"context"
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
)

var (
	genericError = errors.New("generic error")
)

func BuildMembers(size int) []*domain.Member {
	var members []*domain.Member
	for i := 0; i < size; i++ {
		members = append(members, buildMember(domain.NewID()))
	}
	return members
}

func buildMember(id string) *domain.Member {
	return &domain.Member{
		ID: id,
	}
}

func BuildContext() context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{
			ID: "church_id_test",
		},
	})
}
