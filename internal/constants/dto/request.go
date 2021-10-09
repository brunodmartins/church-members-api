package dto

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum/role"
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

//CreateUserRequest for HTTP calls to post user
// swagger:model CreateUserRequest
type CreateUserRequest struct {
	UserName string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email,min=3,max=32"`
	Role     string `json:"role" validate:"required,eq=ADMIN|eq=USER"`
	Password string `json:"password" validate:"required,password"`
}

func (r CreateUserRequest) ToUser() *domain.User {
	return domain.NewUser(r.UserName, r.Email, r.Password, role.From(r.Role))
}
