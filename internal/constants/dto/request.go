package dto

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"time"
)

//PutMemberStatusRequest for HTTP calls to put member status
// swagger:model PutMemberStatusRequest
type PutMemberStatusRequest struct {
	Active *bool     `json:"active" binding:"required"`
	Reason string    `json:"reason" binding:"required"`
	Date   time.Time `json:"date"`
}

//CreateMemberRequest for HTTP calls to post member
// swagger:model CreateMemberRequest
type CreateMemberRequest struct {
	*model.Member
}