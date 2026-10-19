package jobs

import (
	"context"
	"encoding/json"

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
	messages, err := m.queue.ReceiveMessage(ctx, 10)
	if err != nil {
		logrus.Errorf("Error processing queue messages: %v", err)
		return err
	}
	for _, message := range messages {
		var event member.Event
		err := json.Unmarshal([]byte(*message), &event)
		if err != nil {
			logrus.Errorf("Error unmarshalling message: %v", err)
			continue
		}
		logrus.Infof("Processed event: %+v", event)
	}
	return nil
}
