package domain

// Church provides information about a church
type Church struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Language     string `json:"language"`
	Email        string `json:"email"`
}
