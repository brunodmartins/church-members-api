package queue

import (
	"context"
	"errors"
	"testing"

	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestQueue_ReceiveMessage(t *testing.T) {
	t.Run("returns received messages", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSWrapper(ctrl)
		service := New(client)
		ctx := context.Background()
		maxMessages := int64(5)
		messages := []*string{stringPointer("first"), stringPointer("second")}

		client.EXPECT().ReceiveMessage(ctx, maxMessages).Return(messages, nil)

		actual, err := service.ReceiveMessage(ctx, maxMessages)

		assert.NoError(t, err)
		assert.Equal(t, messages, actual)
	})

	t.Run("returns receive error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSWrapper(ctrl)
		service := New(client)
		ctx := context.Background()
		wantErr := errors.New("receive failed")

		client.EXPECT().ReceiveMessage(ctx, int64(1)).Return(nil, wantErr)

		actual, err := service.ReceiveMessage(ctx, 1)

		assert.Nil(t, actual)
		assert.ErrorIs(t, err, wantErr)
	})
}

func TestQueue_SendMessage(t *testing.T) {
	t.Run("sends message successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSWrapper(ctrl)
		service := New(client)
		ctx := context.Background()
		message := []byte("member updated")

		client.EXPECT().SendMessage(ctx, message).Return(nil)

		assert.NoError(t, service.SendMessage(ctx, message))
	})

	t.Run("returns send error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSWrapper(ctrl)
		service := New(client)
		ctx := context.Background()
		message := []byte("member updated")
		wantErr := errors.New("send failed")

		client.EXPECT().SendMessage(ctx, message).Return(wantErr)

		assert.ErrorIs(t, service.SendMessage(ctx, message), wantErr)
	})
}

func stringPointer(value string) *string {
	return &value
}
