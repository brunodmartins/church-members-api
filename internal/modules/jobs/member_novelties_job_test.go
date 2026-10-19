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
		name     string
		messages []*string
		err      error
	}{
		{
			name:     "processes a valid event",
			messages: []*string{stringPointer(`{"memberId":"member-123","eventType":"MEMBER_CREATED"}`)},
		},
		{
			name:     "handles malformed messages and continues",
			messages: []*string{stringPointer("{"), stringPointer(`{"memberId":"member-123","eventType":"MEMBER_CREATED"}`)},
		},
		{
			name: "handles receive errors",
			err:  errors.New("queue unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			queueService := queue.NewMockService(ctrl)
			queueService.EXPECT().ReceiveMessage(ctx, int64(20)).Return(tt.messages, tt.err)

			gotErr := NewMemberNoveltiesJob(queueService).RunJob(ctx)
			if tt.err != nil {
				assert.ErrorIs(t, gotErr, tt.err)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func stringPointer(value string) *string {
	return &value
}
