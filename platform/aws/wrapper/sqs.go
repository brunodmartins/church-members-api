package wrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
)

// SQSAPI defines the interface for the AWS SQS client.
//
//go:generate mockgen -source=./sqs.go -destination=./mock/sqs_mock.go
type SQSAPI interface {
	SendMessage(context.Context, *sqs.SendMessageInput, ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	ReceiveMessage(context.Context, *sqs.ReceiveMessageInput, ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

// SQSWrapper provides a simplified interface for interacting with AWS SQS.
type SQSWrapper interface {
	SendMessage(ctx context.Context, messageBody []byte) error
	ReceiveMessage(ctx context.Context, maxMessages int64) ([]*string, error)
}

type sqsWrapper struct {
	client   SQSAPI
	queueURL string
}

// NewSQSWrapper creates a new instance of SQSWrapper with the given AWS client and queue URL.
func NewSQSWrapper(client SQSAPI, queueURL string) SQSWrapper {
	return &sqsWrapper{
		client:   client,
		queueURL: queueURL,
	}
}

// SendMessage sends a message to the SQS queue.
func (w *sqsWrapper) SendMessage(ctx context.Context, messageBody []byte) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(w.queueURL),
		MessageBody: aws.String(string(messageBody)),
	}
	_, err := w.client.SendMessage(ctx, input)
	if err != nil {
		logrus.Errorf("Error sending message to AWS SQS %s: %v", w.queueURL, err)
	}
	return err
}

// ReceiveMessage retrieves messages from the SQS queue, up to the specified maximum number of messages.
func (w *sqsWrapper) ReceiveMessage(ctx context.Context, maxMessages int64) ([]*string, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(w.queueURL),
		MaxNumberOfMessages: int32(maxMessages),
	}
	output, err := w.client.ReceiveMessage(ctx, input)
	if err != nil {
		logrus.Errorf("Error receiving message from AWS SQS %s: %v", w.queueURL, err)
		return nil, err
	}
	logrus.Infof("Received %d messages from AWS SQS %s", len(output.Messages), w.queueURL)
	messages := make([]*string, len(output.Messages))
	for i, msg := range output.Messages {
		messages[i] = msg.Body
	}
	return messages, nil
}
