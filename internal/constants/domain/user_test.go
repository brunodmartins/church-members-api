package domain

import (
	"context"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/role"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	const password = "123"
	user := NewUser("", "", "", password, role.USER, NotificationPreferences{})
	assert.NotEqual(t, password, string(user.Password))
}

func TestGetChurch(t *testing.T) {
	user := NewUser("", "", "", "", role.USER, NotificationPreferences{})
	user.Church = &Church{ID: "church-id"}
	ctx := context.WithValue(context.TODO(), "user", user)
	assert.Equal(t, "church-id", GetChurchID(ctx))
	assert.NotNil(t, GetChurch(ctx))
}
