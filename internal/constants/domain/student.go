package domain

// Student represents a student entity.
type Student struct {
	ID          string       `json:"id"`
	StudentID   string       `json:"studentId"`
	Name        string       `json:"name"`
	Gender      string       `json:"gender"`
	Member      *Member      `json:"member,omitempty"`
	Participant *Participant `json:"participant,omitempty"`
}
