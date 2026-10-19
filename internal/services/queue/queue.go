package queue

import (
	"context"

	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=queue.go -destination=mock_queue.go -package=queue
type Service interface {
	ReceiveMessage(ctx context.Context, maxMessages int64) ([]*string, error)
	SendMessage(ctx context.Context, messageBody []byte) error
}

type defaultQueue struct {
	client wrapper.SQSWrapper
}

func New(client wrapper.SQSWrapper) Service {
	return &defaultQueue{
		client: client,
	}
}

// ReceiveMessage retrieves messages from the SQS queue.
func (q *defaultQueue) ReceiveMessage(ctx context.Context, maxMessages int64) ([]*string, error) {
	logrus.Info("Polling SQS for new messages")
	return q.client.ReceiveMessage(ctx, maxMessages)
}

// SendMessage sends a message to the SQS queue.
func (q *defaultQueue) SendMessage(ctx context.Context, messageBody []byte) error {
	logrus.Info("Sending message to SQS")
	return q.client.SendMessage(ctx, messageBody)
}
