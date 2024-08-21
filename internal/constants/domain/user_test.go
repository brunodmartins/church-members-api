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
	assert.False(t, user.ConfirmedEmail)
}

func TestGetChurchID(t *testing.T) {
	t.Run("Given a valid context with 'user', when request the church id from the context, then return it", func(t *testing.T) {
		user := NewUser("", "", "", "", role.USER, NotificationPreferences{})
		user.Church = &Church{ID: "church-id"}
		ctx := context.WithValue(context.TODO(), "user", user)
		assert.Equal(t, "church-id", GetChurchID(ctx))
		assert.NotNil(t, GetChurch(ctx))
	})
	t.Run("Given a valid context with 'church', when request the church id from the context, then return it", func(t *testing.T) {
		church := &Church{ID: "church-id"}
		ctx := context.WithValue(context.TODO(), "church", church)
		assert.Equal(t, "church-id", GetChurchID(ctx))
		assert.NotNil(t, GetChurch(ctx))
	})
	t.Run("Given a valid context with only 'church_id', when request the church id from the context, then return it", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), "church_id", "church-id")
		assert.Equal(t, "church-id", GetChurchID(ctx))
		assert.Nil(t, GetChurch(ctx))
	})

}
