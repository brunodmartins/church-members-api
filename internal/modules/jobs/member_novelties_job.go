package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/brunodmartins/church-members-api/internal/constants/enum/reportType"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/modules/report"
	"github.com/brunodmartins/church-members-api/internal/services/queue"
	"github.com/sirupsen/logrus"
)

type memberNoveltiesJob struct {
	queue   queue.Service
	reports report.Service
}

func NewMemberNoveltiesJob(queue queue.Service, reports report.Service) Job {
	return &memberNoveltiesJob{
		queue:   queue,
		reports: reports,
	}
}

func (m *memberNoveltiesJob) RunJob(ctx context.Context) error {
	messages, err := m.queue.ReceiveMessages(ctx, 10)
	if err != nil {
		logrus.Errorf("Error processing queue messages: %v", err)
		return err
	}

	var processingErrors []error
	shouldRegenerateReports := false
	var processedMessages []queue.Message
	for _, message := range messages {
		var event member.Event
		if err := json.Unmarshal([]byte(message.Body), &event); err != nil {
			logrus.Errorf("Error unmarshalling message: %v", err)
			continue
		}
		logrus.Infof("Processed event: %+v", event)
		switch event.EventType {
		case member.EventMemberCreated, member.EventMemberUpdated, member.EventMemberRetired:
			shouldRegenerateReports = true
		}
		processedMessages = append(processedMessages, message)
	}
	if shouldRegenerateReports {
		if err := m.generateReports(ctx); err != nil {
			return err
		}
	}
	for _, message := range processedMessages {
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

func (m *memberNoveltiesJob) generateReports(ctx context.Context) error {
	var generationErrors []error
	for _, reportName := range reportType.ReportsTypes {
		if err := m.reports.GenerateReport(ctx, reportName); err != nil {
			logrus.WithError(err).Errorf("Error generating report %s", reportName)
			generationErrors = append(generationErrors, fmt.Errorf("generate %s report: %w", reportName, err))
		}
	}
	return errors.Join(generationErrors...)
}
