package dto

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"time"
)

//PutMemberStatusRequest for HTTP calls to put member status
// swagger:model PutMemberStatusRequest
type PutMemberStatusRequest struct {
	Active *bool     `json:"active" validate:"required"`
	Reason string    `json:"reason" validate:"required"`
	Date   time.Time `json:"date"`
}

//CreateMemberRequest for HTTP calls to post member
// swagger:model CreateMemberRequest
type CreateMemberRequest struct {
	*domain.Member
}