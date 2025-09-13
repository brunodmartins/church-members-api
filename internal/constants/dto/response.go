package dto

import (
	"time"
)

// ErrorResponse for HTTP error responses
// swagger:model ErrorResponse
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// MessageResponse for HTTP success responses
// swagger:model MessageResponse
type MessageResponse struct {
	Message string `json:"message"`
}

// CreateMemberResponse for HTTP create member responses
// swagger:model CreateMemberResponse
type CreateMemberResponse struct {
	ID string `json:"id"`
}

// GetMemberResponse for HTTP get member responses
// swagger:model GetMemberResponse
type GetMemberResponse struct {
	ID             string             `json:"id"`
	Active         bool               `json:"active"`
	Classification string             `json:"classification"`
	Person         *GetPersonResponse `json:"person"`
}

// GetPersonResponse for HTTP get person response
// swagger:model GetPersonResponse
type GetPersonResponse struct {
	FirstName     string              `json:"firstName,omitempty"`
	LastName      string              `json:"lastName,omitempty"`
	FullName      string              `json:"fullName,omitempty"`
	Gender        string              `json:"gender,omitempty"`
	Age           int                 `json:"age,omitempty"`
	BirthDate     time.Time           `json:"birthDate,omitempty"`
	MarriageDate  *time.Time          `json:"marriageDate,omitempty"`
	SpousesName   string              `json:"spousesName,omitempty"`
	MaritalStatus string              `json:"maritalStatus,omitempty"`
	Contact       *GetContactResponse `json:"contact,omitempty"`
	Address       *GetAddressResponse `json:"address,omitempty"`
}

// GetContactResponse for HTTP get contact response
// swagger:model GetContactResponse
type GetContactResponse struct {
	Cellphone string `json:"cellphone,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
}

// GetAddressResponse for HTTP get address response
// swagger:model GetAddressResponse
type GetAddressResponse struct {
	ZipCode  string `json:"zipCode,omitempty"`
	State    string `json:"state,omitempty"`
	City     string `json:"city,omitempty"`
	Address  string `json:"address,omitempty"`
	District string `json:"district,omitempty"`
	Number   int    `json:"number,omitempty"`
	Full     string `json:"full,omitempty"`
}

// GetTokenResponse for HTTP get token responses
// swagger:model GetTokenResponse
type GetTokenResponse struct {
	Token string `json:"token"`
}
