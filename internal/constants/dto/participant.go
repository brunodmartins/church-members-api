package dto

import (
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

// CreateParticipantRequest for HTTP calls
type CreateParticipantRequest struct {
	Name        string `json:"name" validate:"required"`
	BirthDate   Date   `json:"birthDate"`
	CellPhone   string `json:"cellPhone"`
	Filiation   string `json:"filiation"`
	Observation string `json:"observation"`
}

func (r *CreateParticipantRequest) ToParticipant() *domain.Participant {
	return &domain.Participant{
		Name:        r.Name,
		BirthDate:   *ToTime(&r.BirthDate),
		CellPhone:   r.CellPhone,
		Filiation:   r.Filiation,
		Observation: r.Observation,
	}
}

// GetParticipantResponse for HTTP responses
type GetParticipantResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BirthDate   *Date  `json:"birthDate,omitempty"`
	CellPhone   string `json:"cellPhone,omitempty"`
	Filiation   string `json:"filiation,omitempty"`
	Observation string `json:"observation,omitempty"`
}

func NewGetParticipantResponse(p *domain.Participant) *GetParticipantResponse {
	if p == nil {
		return nil
	}
	var bd *Date
	d := Date{p.BirthDate}
	bd = &d
	return &GetParticipantResponse{
		ID:          p.ID,
		Name:        p.Name,
		BirthDate:   bd,
		CellPhone:   p.CellPhone,
		Filiation:   p.Filiation,
		Observation: p.Observation,
	}
}
