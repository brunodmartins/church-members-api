package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	const password = "123"
	user := NewUser("", "", "", password, NotificationPreferences{})
	assert.NotEqual(t, password, string(user.Password))
	assert.False(t, user.ConfirmedEmail)
}

func TestGetChurchID(t *testing.T) {
	t.Run("Given a valid context with 'user', when request the church id from the context, then return it", func(t *testing.T) {
		user := NewUser("", "", "", "", NotificationPreferences{})
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
