package participant

import (
	"context"
	"errors"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

var genericError = errors.New("generic error")

func BuildParticipants(size int) []*domain.Participant {
	var participants []*domain.Participant
	for i := 0; i < size; i++ {
		participants = append(participants, buildParticipant(domain.NewID()))
	}
	return participants
}

func buildParticipant(id string) *domain.Participant {
	now := time.Now()
	return &domain.Participant{
		ID:          id,
		ChurchID:    "church_id_test",
		Name:        "First Last",
		BirthDate:   now,
		CellPhone:   "99999999",
		Filiation:   "Filiation",
		Observation: "Obs",
	}
}

func BuildContext() context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{ID: "church_id_test"},
	})
}
