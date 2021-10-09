package domain

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"github.com/BrunoDM2943/church-members-api/platform/crypto"
)

type User struct {
	ID       string    `json:"id"`
	UserName string    `json:"username"`
	Email    string    `json:"email"`
	Role     enum.Role `json:"role"`
	Password []byte    `json:"-"`
}

func NewUser(userName, email, password string, role enum.Role) *User {
	return &User{UserName: userName, Email: email, Password: crypto.EncryptPassword(password), Role: role}
}

