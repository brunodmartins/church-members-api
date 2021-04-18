package dto

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/graphql-go/graphql/gqlerrors"
)

//ErrorResponse for HTTP error responses
// swagger:model ErrorResponse
type ErrorResponse struct {
	Message string
	Error   error
}

//GraphQLErrorResponse for HTTP error responses
// swagger:model GraphQLErrorResponse
type GraphQLErrorResponse struct {
	Errors []gqlerrors.FormattedError
}

//CreateMemberResponse for HTTP create member responses
// swagger:model CreateMemberResponse
type CreateMemberResponse struct {
	ID		string
}

//GetMemberResponse for HTTP get member responses
// swagger:model GetMemberResponse
type GetMemberResponse struct {
	*model.Member
}