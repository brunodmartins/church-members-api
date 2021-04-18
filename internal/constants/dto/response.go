package dto

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/graphql-go/graphql/gqlerrors"
)

type ErrorResponse struct {
	Message string
	Error   error
}

type GraphQLErrorResponse struct {
	Errors []gqlerrors.FormattedError
}

type CreateMemberResponse struct {
	ID		string
}

type GetMemberResponse struct {
	*model.Member
}