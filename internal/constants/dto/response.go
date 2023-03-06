package dto

import (
	"github.com/graphql-go/graphql/gqlerrors"
	"time"
)

// ErrorResponse for HTTP error responses
// swagger:model ErrorResponse
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// GraphQLErrorResponse for HTTP error responses
// swagger:model GraphQLErrorResponse
type GraphQLErrorResponse struct {
	Errors []gqlerrors.FormattedError `json:"errors"`
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
	FirstName    string              `json:"firstName"`
	LastName     string              `json:"lastName"`
	FullName     string              `json:"fullName"`
	Gender       string              `json:"gender"`
	Age          int                 `json:"age"`
	BirthDate    time.Time           `json:"birthDate"`
	MarriageDate *time.Time          `json:"marriageDate"`
	SpousesName  string              `json:"spousesName"`
	Contact      *GetContactResponse `json:"contact"`
	Address      *GetAddressResponse `json:"address"`
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
	ZipCode  string `json:"zipCode"`
	State    string `json:"state"`
	City     string `json:"city"`
	Address  string `json:"address"`
	District string `json:"district"`
	Number   int    `json:"number"`
	Full     string `json:"full"`
}

// GetTokenResponse for HTTP get token responses
// swagger:model GetTokenResponse
type GetTokenResponse struct {
	Token string `json:"token"`
}
