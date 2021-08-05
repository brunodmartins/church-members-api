package dto

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/graphql-go/graphql/gqlerrors"
)

//ErrorResponse for HTTP error responses
// swagger:model ErrorResponse
type ErrorResponse struct {
	Message string `json:"message"`
	Error   error `json:"error"`
}

//GraphQLErrorResponse for HTTP error responses
// swagger:model GraphQLErrorResponse
type GraphQLErrorResponse struct {
	Errors []gqlerrors.FormattedError `json:"errors"`
}

//CreateMemberResponse for HTTP create member responses
// swagger:model CreateMemberResponse
type CreateMemberResponse struct {
	ID		string `json:"id"`
}

//GetMemberResponse for HTTP get member responses
// swagger:model GetMemberResponse
type GetMemberResponse struct {
	*domain.Member
}