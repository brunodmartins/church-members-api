package queue

import (
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
		maxMessages := int64(5)
		messages := []*string{stringPointer("first"), stringPointer("second")}

		client.EXPECT().ReceiveMessage(maxMessages).Return(messages, nil)

		actual, err := service.ReceiveMessage(maxMessages)

		assert.NoError(t, err)
		assert.Equal(t, messages, actual)
	})

	t.Run("returns receive error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSWrapper(ctrl)
		service := New(client)
		wantErr := errors.New("receive failed")

		client.EXPECT().ReceiveMessage(int64(1)).Return(nil, wantErr)

		actual, err := service.ReceiveMessage(1)

		assert.Nil(t, actual)
		assert.ErrorIs(t, err, wantErr)
	})
}

func TestQueue_SendMessage(t *testing.T) {
	t.Run("sends message successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSWrapper(ctrl)
		service := New(client)
		message := []byte("member updated")

		client.EXPECT().SendMessage(message).Return(nil)

		assert.NoError(t, service.SendMessage(message))
	})

	t.Run("returns send error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSWrapper(ctrl)
		service := New(client)
		message := []byte("member updated")
		wantErr := errors.New("send failed")

		client.EXPECT().SendMessage(message).Return(wantErr)

		assert.ErrorIs(t, service.SendMessage(message), wantErr)
	})
}

func stringPointer(value string) *string {
	return &value
}
