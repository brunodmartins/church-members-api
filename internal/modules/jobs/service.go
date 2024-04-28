package jobs

import (
	"context"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

// Job exposing jobs operations
//
//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Job interface {
	RunJob(ctx context.Context) error
}

func fmtDate(date time.Time) string {
	return date.Format("02-Jan")
}

func getEmail(user *domain.User) string {
	return user.Email
}

func getPhone(user *domain.User) string {
	return user.Phone
}

func mapToSlice(users []*domain.User, mapper func(user *domain.User) string) []string {
	var result []string
	for _, user := range users {
		result = append(result, mapper(user))
	}
	return result
}
