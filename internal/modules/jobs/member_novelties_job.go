package jobs

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/services/queue"
	"github.com/sirupsen/logrus"
)

type memberNoveltiesJob struct {
	queue queue.Service
}

func NewMemberNoveltiesJob(queue queue.Service) Job {
	return &memberNoveltiesJob{
		queue: queue,
	}
}

func (m *memberNoveltiesJob) RunJob(ctx context.Context) error {
	messages, err := m.queue.ReceiveMessages(ctx, 10)
	if err != nil {
		logrus.Errorf("Error processing queue messages: %v", err)
		return err
	}

	var processingErrors []error
	for _, message := range messages {
		var event member.Event
		if err := json.Unmarshal([]byte(message.Body), &event); err != nil {
			logrus.Errorf("Error unmarshalling message: %v", err)
			continue
		}
		logrus.Infof("Processed event: %+v", event)
		if err := m.deleteMessage(ctx, message); err != nil {
			processingErrors = append(processingErrors, err)
		}
	}
	return errors.Join(processingErrors...)
}

func (m *memberNoveltiesJob) deleteMessage(ctx context.Context, message queue.Message) error {
	if err := m.queue.DeleteMessage(ctx, message.ReceiptHandle); err != nil {
		logrus.WithError(err).Error("Error deleting queue message")
		return err
	}
	return nil
}
