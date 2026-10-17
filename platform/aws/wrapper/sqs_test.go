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

		client.EXPECT().SendMessage(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, input *sqs.SendMessageInput, _ ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			assert.Equal(t, queueURL, *input.QueueUrl)
			assert.Equal(t, string(messageBody), *input.MessageBody)
			return &sqs.SendMessageOutput{}, nil
		})

		assert.NoError(t, service.SendMessage(messageBody))
	})

	t.Run("returns send error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSAPI(ctrl)
		service := NewSQSWrapper(client, queueURL)
		wantErr := errors.New("send failed")

		client.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil, wantErr)

		assert.ErrorIs(t, service.SendMessage(messageBody), wantErr)
	})
}

func TestSQSWrapper_ReceiveMessage(t *testing.T) {
	const queueURL = "https://sqs.us-east-1.amazonaws.com/123456789012/members"
	maxMessages := int64(3)

	t.Run("returns received message bodies", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSAPI(ctrl)
		service := NewSQSWrapper(client, queueURL)
		output := &sqs.ReceiveMessageOutput{Messages: []types.Message{
			{Body: aws.String("first")},
			{Body: aws.String("second")},
		}}

		client.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, input *sqs.ReceiveMessageInput, _ ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
			assert.Equal(t, queueURL, *input.QueueUrl)
			assert.Equal(t, int32(maxMessages), input.MaxNumberOfMessages)
			return output, nil
		})

		messages, err := service.ReceiveMessage(maxMessages)

		assert.NoError(t, err)
		assert.Equal(t, []*string{aws.String("first"), aws.String("second")}, messages)
	})

	t.Run("returns receive error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		client := mock_wrapper.NewMockSQSAPI(ctrl)
		service := NewSQSWrapper(client, queueURL)
		wantErr := errors.New("receive failed")

		client.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).Return(nil, wantErr)

		messages, err := service.ReceiveMessage(maxMessages)

		assert.Nil(t, messages)
		assert.ErrorIs(t, err, wantErr)
	})
}
