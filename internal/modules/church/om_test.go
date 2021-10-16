package church

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
)

var genericError = errors.New("generic error")

func buildChurch(id string) *domain.Church {
	return &domain.Church{
		ID:           id,
		Name:         "test church",
		Abbreviation: "tc",
		Language:     "pt-br",
	}
}
