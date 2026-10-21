package jobs

import (
	"context"
	"errors"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/services/queue"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMemberNoveltiesJob_RunJob(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		messages      []queue.Message
		receiveErr    error
		expectDeleted []string
	}{
		{
			name: "processes and acknowledges received events",
			messages: []queue.Message{
				{Body: `{"churchId":"church-1","memberId":"member-123","eventType":"MEMBER_CREATED"}`, ReceiptHandle: "matching"},
			},
			expectDeleted: []string{"matching"},
		},
		{
			name: "leaves malformed messages in the queue",
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

			gotErr := NewMemberNoveltiesJob(queueService).RunJob(ctx)
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
		{Body: `{"churchId":"church-1","memberId":"member-123","eventType":"MEMBER_CREATED"}`, ReceiptHandle: "receipt"},
	}, nil)
	queueService.EXPECT().DeleteMessage(ctx, "receipt").Return(deleteErr)

	assert.ErrorIs(t, NewMemberNoveltiesJob(queueService).RunJob(ctx), deleteErr)
}
