package queue

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestQueue_ReceiveMessages(t *testing.T) {
	t.Run("returns only messages for the church in the context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSWrapper(ctrl)
		service := New(client)
		ctx := context.WithValue(context.Background(), "church", &domain.Church{ID: "church-1"})
		maxMessages := int64(5)
		messages := []types.Message{
			{Body: aws.String(`{"churchId":"church-1","eventType":"MEMBER_CREATED"}`), ReceiptHandle: aws.String("matching-receipt")},
			{Body: aws.String(`{"churchId":"church-2","eventType":"MEMBER_CREATED"}`), ReceiptHandle: aws.String("foreign-receipt")},
			{Body: aws.String("{"), ReceiptHandle: aws.String("malformed-receipt")},
		}
		client.EXPECT().ReceiveMessages(ctx, maxMessages).Return(messages, nil)

		actual, err := service.ReceiveMessages(ctx, maxMessages)

		assert.NoError(t, err)
		assert.Equal(t, []Message{
			{Body: `{"churchId":"church-1","eventType":"MEMBER_CREATED"}`, ReceiptHandle: "matching-receipt"},
		}, actual)
	})

	t.Run("returns receive error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSWrapper(ctrl)
		service := New(client)
		ctx := context.Background()
		wantErr := errors.New("receive failed")

		client.EXPECT().ReceiveMessages(ctx, int64(1)).Return(nil, wantErr)

		actual, err := service.ReceiveMessages(ctx, 1)

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

func TestQueue_DeleteMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock_wrapper.NewMockSQSWrapper(ctrl)
	service := New(client)
	ctx := context.Background()

	client.EXPECT().DeleteMessage(ctx, "receipt-handle").Return(nil)

	assert.NoError(t, service.DeleteMessage(ctx, "receipt-handle"))
}
