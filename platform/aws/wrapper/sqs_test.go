package wrapper

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSQSWrapper_SendMessage(t *testing.T) {
	const queueURL = "https://sqs.us-east-1.amazonaws.com/123456789012/members"
	messageBody := []byte("member updated")

	t.Run("sends message successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSAPI(ctrl)
		service := NewSQSWrapper(client, queueURL)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		client.EXPECT().SendMessage(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, input *sqs.SendMessageInput, _ ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			assert.Equal(t, queueURL, *input.QueueUrl)
			assert.Equal(t, string(messageBody), *input.MessageBody)
			return &sqs.SendMessageOutput{}, nil
		})

		assert.NoError(t, service.SendMessage(ctx, messageBody))
	})

	t.Run("returns send error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSAPI(ctrl)
		service := NewSQSWrapper(client, queueURL)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wantErr := errors.New("send failed")

		client.EXPECT().SendMessage(ctx, gomock.Any()).Return(nil, wantErr)

		assert.ErrorIs(t, service.SendMessage(ctx, messageBody), wantErr)
	})
}

func TestSQSWrapper_ReceiveMessages(t *testing.T) {
	const queueURL = "https://sqs.us-east-1.amazonaws.com/123456789012/members"
	maxMessages := int64(3)

	t.Run("returns received message bodies", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSAPI(ctrl)
		service := NewSQSWrapper(client, queueURL)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		output := &sqs.ReceiveMessageOutput{Messages: []types.Message{
			{Body: aws.String("first"), ReceiptHandle: aws.String("first-receipt")},
			{Body: aws.String("second"), ReceiptHandle: aws.String("second-receipt")},
		}}

		client.EXPECT().ReceiveMessage(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, input *sqs.ReceiveMessageInput, _ ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
			assert.Equal(t, queueURL, *input.QueueUrl)
			assert.Equal(t, int32(maxMessages), input.MaxNumberOfMessages)
			return output, nil
		})

		messages, err := service.ReceiveMessages(ctx, maxMessages)

		assert.NoError(t, err)
		assert.Equal(t, output.Messages, messages)
	})

	t.Run("returns receive error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSAPI(ctrl)
		service := NewSQSWrapper(client, queueURL)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wantErr := errors.New("receive failed")

		client.EXPECT().ReceiveMessage(ctx, gomock.Any()).Return(nil, wantErr)

		messages, err := service.ReceiveMessages(ctx, maxMessages)

		assert.Nil(t, messages)
		assert.ErrorIs(t, err, wantErr)
	})
}

func TestSQSWrapper_DeleteMessage(t *testing.T) {
	const queueURL = "https://sqs.us-east-1.amazonaws.com/123456789012/members"
	ctrl := gomock.NewController(t)
	client := mock_wrapper.NewMockSQSAPI(ctrl)
	service := NewSQSWrapper(client, queueURL)
	ctx := context.Background()

	client.EXPECT().DeleteMessage(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, input *sqs.DeleteMessageInput, _ ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
		assert.Equal(t, queueURL, aws.ToString(input.QueueUrl))
		assert.Equal(t, "receipt-handle", aws.ToString(input.ReceiptHandle))
		return &sqs.DeleteMessageOutput{}, nil
	})

	assert.NoError(t, service.DeleteMessage(ctx, "receipt-handle"))
}
