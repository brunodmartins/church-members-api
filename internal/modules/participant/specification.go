package participant

import "github.com/brunodmartins/church-members-api/internal/constants/domain"

type Specification func(participant *domain.Participant) bool

func applySpecifications(participants []*domain.Participant, specification []Specification) []*domain.Participant {
	var filtered []*domain.Participant
	for _, participant := range participants {
		allSpecTrue := true
		for _, spec := range specification {
			allSpecTrue = allSpecTrue && spec(participant)
		}
		if allSpecTrue {
			filtered = append(filtered, participant)
		}
	}
	return filtered
}
