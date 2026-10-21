package jobs

import (
	"context"
	"errors"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/enum/reportType"
	mock_report "github.com/brunodmartins/church-members-api/internal/modules/report/mock"
	"github.com/brunodmartins/church-members-api/internal/services/queue"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMemberNoveltiesJob_RunJob(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name              string
		messages          []queue.Message
		receiveErr        error
		expectDeleted     []string
		regenerateReports bool
	}{
		{
			name: "acknowledges created member events and regenerates reports",
			messages: []queue.Message{
				{Body: `{"churchId":"church-1","memberId":"member-123","eventType":"MEMBER_CREATED"}`, ReceiptHandle: "created"},
			},
			expectDeleted:     []string{"created"},
			regenerateReports: true,
		},
		{
			name: "leaves malformed messages queued and regenerates after member update",
			messages: []queue.Message{
				{Body: "{", ReceiptHandle: "malformed"},
				{Body: `{"memberId":"member-123","eventType":"MEMBER_UPDATED"}`, ReceiptHandle: "updated"},
			},
			expectDeleted:     []string{"updated"},
			regenerateReports: true,
		},
		{
			name: "acknowledges retired member events and regenerates reports",
			messages: []queue.Message{
				{Body: `{"memberId":"member-123","eventType":"MEMBER_RETIRED"}`, ReceiptHandle: "retired"},
			},
			expectDeleted:     []string{"retired"},
			regenerateReports: true,
		},
		{
			name: "acknowledges unrelated events without regenerating reports",
			messages: []queue.Message{
				{Body: `{"memberId":"member-123","eventType":"OTHER"}`, ReceiptHandle: "other"},
			},
			expectDeleted: []string{"other"},
		},
		{
			name: "leaves malformed messages queued",
			messages: []queue.Message{
				{Body: "{", ReceiptHandle: "malformed"},
			},
		},
		{
			name:       "returns receive errors",
			receiveErr: errors.New("queue unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			queueService := queue.NewMockService(ctrl)
			queueService.EXPECT().ReceiveMessages(ctx, int64(10)).Return(tt.messages, tt.receiveErr)
			for _, receiptHandle := range tt.expectDeleted {
				queueService.EXPECT().DeleteMessage(ctx, receiptHandle).Return(nil)
			}

			reportService := mock_report.NewMockService(ctrl)
			if tt.regenerateReports {
				for _, reportName := range reportType.ReportsTypes {
					reportService.EXPECT().GenerateReport(ctx, reportName).Return(nil)
				}
			}

			gotErr := NewMemberNoveltiesJob(queueService, reportService).RunJob(ctx)
			if tt.receiveErr != nil {
				assert.ErrorIs(t, gotErr, tt.receiveErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestMemberNoveltiesJob_RunJob_ReturnsDeleteErrors(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	queueService := queue.NewMockService(ctrl)
	deleteErr := errors.New("delete failed")
	queueService.EXPECT().ReceiveMessages(ctx, int64(10)).Return([]queue.Message{
		{Body: `{"churchId":"church-1","memberId":"member-123","eventType":"OTHER"}`, ReceiptHandle: "receipt"},
	}, nil)
	queueService.EXPECT().DeleteMessage(ctx, "receipt").Return(deleteErr)

	assert.ErrorIs(t, NewMemberNoveltiesJob(queueService, mock_report.NewMockService(ctrl)).RunJob(ctx), deleteErr)
}

func TestMemberNoveltiesJob_RunJob_DoesNotDeleteMessagesWhenReportGenerationFails(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	queueService := queue.NewMockService(ctrl)
	queueService.EXPECT().ReceiveMessages(ctx, int64(10)).Return([]queue.Message{
		{Body: `{"memberId":"member-123","eventType":"MEMBER_CREATED"}`, ReceiptHandle: "receipt"},
	}, nil)

	reportErr := errors.New("report generation failed")
	reportService := mock_report.NewMockService(ctrl)
	for index, reportName := range reportType.ReportsTypes {
		generationErr := error(nil)
		if index == 0 {
			generationErr = reportErr
		}
		reportService.EXPECT().GenerateReport(ctx, reportName).Return(generationErr)
	}

	err := NewMemberNoveltiesJob(queueService, reportService).RunJob(ctx)
	assert.ErrorIs(t, err, reportErr)
}
