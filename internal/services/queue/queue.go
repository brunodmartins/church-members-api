package queue

import (
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=queue.go -destination=mock_queue.go -package=queue
type Service interface {
	ReceiveMessage(maxMessages int64) ([]*string, error)
	SendMessage(messageBody []byte) error
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
func (q *defaultQueue) ReceiveMessage(maxMessages int64) ([]*string, error) {
	logrus.Info("Polling SQS for new messages")
	return q.client.ReceiveMessage(maxMessages)
}

// SendMessage sends a message to the SQS queue.
func (q *defaultQueue) SendMessage(messageBody []byte) error {
	logrus.Info("Sending message to SQS")
	return q.client.SendMessage(messageBody)
}
