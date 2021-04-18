package dto

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"time"
)

type PutMemberStatusRequest struct {
	Active *bool     `json:"active" binding:"required"`
	Reason string    `json:"reason" binding:"required"`
	Date   time.Time `json:"date"`
}

type CreateMemberRequest struct {
	*model.Member
}