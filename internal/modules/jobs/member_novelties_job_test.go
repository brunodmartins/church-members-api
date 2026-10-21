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
		messages          []*string
		err               error
		regenerateReports bool
	}{
		{
			name:              "regenerates reports for a created member",
			messages:          []*string{stringPointer(`{"memberId":"member-123","eventType":"MEMBER_CREATED"}`)},
			regenerateReports: true,
		},
		{
			name:              "regenerates reports after malformed messages and continues",
			messages:          []*string{stringPointer("{"), stringPointer(`{"memberId":"member-123","eventType":"MEMBER_UPDATED"}`)},
			regenerateReports: true,
		},
		{
			name:              "regenerates reports for a retired member",
			messages:          []*string{stringPointer(`{"memberId":"member-123","eventType":"MEMBER_RETIRED"}`)},
			regenerateReports: true,
		},
		{
			name:     "does not regenerate reports for unrelated events",
			messages: []*string{stringPointer(`{"memberId":"member-123","eventType":"OTHER"}`)},
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
			reportService := mock_report.NewMockService(ctrl)
			queueService.EXPECT().ReceiveMessage(ctx, int64(10)).Return(tt.messages, tt.err)
			if tt.regenerateReports {
				for _, reportName := range reportType.ReportsTypes {
					reportService.EXPECT().GenerateReport(ctx, reportName).Return(nil)
				}
			}

			gotErr := NewMemberNoveltiesJob(queueService, reportService).RunJob(ctx)
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
