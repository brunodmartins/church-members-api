package domain

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum/role"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	const password = "123"
	user := NewUser("", "", "", password, role.USER, NotificationPreferences{})
	assert.NotEqual(t, password, string(user.Password))
}
