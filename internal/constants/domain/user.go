package domain

import (
	"context"

	"github.com/brunodmartins/church-members-api/internal/constants/enum"
	"github.com/brunodmartins/church-members-api/platform/crypto"
)

type User struct {
	ID          string                  `json:"id"`
	ChurchID    string                  `json:"church_id"`
	UserName    string                  `json:"username"`
	Email       string                  `json:"email"`
	Role        enum.Role               `json:"role"`
	Phone       string                  `json:"phone"`
	Preferences NotificationPreferences `json:"-"`
	Password    []byte                  `json:"-"`
	Church      *Church                 `json:"-"`
}

type NotificationPreferences struct {
	SendDailySMS    bool `json:"send_daily_sms"`
	SendWeeklyEmail bool `json:"send_weekly_email"`
}

func NewUser(userName, email, password, phone string, role enum.Role, preferences NotificationPreferences) *User {
	return &User{
		UserName:    userName,
		Email:       email,
		Phone:       phone,
		Password:    crypto.EncryptPassword(password),
		Role:        role,
		Preferences: preferences,
	}
}

func GetChurchID(ctx context.Context) string {
	return GetChurch(ctx).ID
}

func GetChurch(ctx context.Context) *Church {
	return ctx.Value("user").(*User).Church
}
