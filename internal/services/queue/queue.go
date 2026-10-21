package queue

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=queue.go -destination=mock_queue.go -package=queue
type Service interface {
	ReceiveMessages(ctx context.Context, maxMessages int64) ([]Message, error)
	SendMessage(ctx context.Context, messageBody []byte) error
	DeleteMessage(ctx context.Context, receiptHandle string) error
}

type Message struct {
	Body          string
	ReceiptHandle string
}

type messageContext struct {
	ChurchID string `json:"churchId"`
}

type defaultQueue struct {
	client wrapper.SQSWrapper
}

func New(client wrapper.SQSWrapper) Service {
	return &defaultQueue{
		client: client,
	}
}

// ReceiveMessages retrieves messages for the church in the context.
func (q *defaultQueue) ReceiveMessages(ctx context.Context, maxMessages int64) ([]Message, error) {
	logrus.Info("Polling SQS for new messages")
	messages, err := q.client.ReceiveMessages(ctx, maxMessages)
	if err != nil {
		return nil, err
	}
	churchID := domain.GetChurchID(ctx)
	result := make([]Message, 0, len(messages))
	for _, message := range messages {
		body := aws.ToString(message.Body)
		var messageContext messageContext
		if err := json.Unmarshal([]byte(body), &messageContext); err != nil {
			logrus.WithError(err).Error("Skipping malformed queue message")
			continue
		}
		if messageContext.ChurchID != churchID {
			continue
		}
		result = append(result, Message{
			Body:          body,
			ReceiptHandle: aws.ToString(message.ReceiptHandle),
		})
	}
	return result, nil
}

// SendMessage sends a message to the SQS queue.
func (q *defaultQueue) SendMessage(ctx context.Context, messageBody []byte) error {
	logrus.Info("Sending message to SQS")
	return q.client.SendMessage(ctx, messageBody)
}

// DeleteMessage acknowledges a successfully processed SQS message.
func (q *defaultQueue) DeleteMessage(ctx context.Context, receiptHandle string) error {
	return q.client.DeleteMessage(ctx, receiptHandle)
}
