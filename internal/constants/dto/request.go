package dto

import (
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/role"
	"time"
)

// RetireMemberRequest for HTTP calls to put member status
// swagger:model RetireMemberRequest
type RetireMemberRequest struct {
	Reason     string    `json:"reason" validate:"required"`
	RetireDate time.Time `json:"date"`
}

// CreateMemberRequest for HTTP calls to post member
// swagger:model CreateMemberRequest
type CreateMemberRequest struct {
	*domain.Member
}

// CreateUserRequest for HTTP calls to post user
// swagger:model CreateUserRequest
type CreateUserRequest struct {
	UserName                       string `json:"username" validate:"required,min=3,max=32"`
	Email                          string `json:"email" validate:"required,email,min=3,max=32"`
	Role                           string `json:"role" validate:"required,eq=ADMIN|eq=USER"`
	Password                       string `json:"password" validate:"required,password"`
	Phone                          string `json:"phone" validate:"required"`
	domain.NotificationPreferences `json:"preferences"`
}

func (r CreateUserRequest) ToUser() *domain.User {
	return domain.NewUser(r.UserName, r.Email, r.Password, r.Phone, role.From(r.Role), r.NotificationPreferences)
}
