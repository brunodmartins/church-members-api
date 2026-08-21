package report

import (
	"context"
	"testing"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/reportType"
	"github.com/brunodmartins/church-members-api/internal/services/storage"
	mock_storage "github.com/brunodmartins/church-members-api/internal/services/storage/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReportService_getReportsCreationDate(t *testing.T) {
	const lastModified = "2026-08-21T10:20:30Z"
	parsedLastModified, _ := time.Parse(time.RFC3339, lastModified)

	tests := []struct {
		name       string
		metadata   storage.FileMetadata
		storageErr error
		expected   string
	}{
		{
			name:     "formats last modified date",
			metadata: storage.FileMetadata{LastModified: &parsedLastModified},
			expected: "21/08/2026 10:20:30",
		},
		{
			name:     "returns fallback when date is missing",
			metadata: storage.FileMetadata{},
			expected: "Error getting last modified date",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storageService := mock_storage.NewMockService(ctrl)
			service := &reportService{storageService: storageService}
			ctx := context.WithValue(context.Background(), "user", &domain.User{
				Church: &domain.Church{ID: "church_id_test"},
			})

			for range reportType.ReportsTypes {
				storageService.EXPECT().GetFileMetadata(gomock.Eq(ctx), gomock.Any()).Return(test.metadata, test.storageErr)
			}

			result := service.getReportsCreationDate(ctx)

			assert.Len(t, result, len(reportType.ReportsTypes))
			for _, reportTypeName := range reportType.ReportsTypes {
				assert.Equal(t, test.expected, result[reportTypeName])
			}
		})
	}
}
