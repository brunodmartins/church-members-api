package domain

import "time"

// Participant represents a participant entity
type Participant struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ChurchID    string    `json:"churchId"`
	Gender      string    `json:"gender,omitempty"`
	BirthDate   time.Time `json:"birthDate,omitempty"`
	CellPhone   string    `json:"cellPhone,omitempty"`
	Filiation   string    `json:"filiation,omitempty"`
	Observation string    `json:"observation,omitempty"`
}
