package jobs

import (
	"context"
	"testing"

	mock_church "github.com/brunodmartins/church-members-api/internal/modules/church/mock"
	mock_jobs "github.com/brunodmartins/church-members-api/internal/modules/jobs/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChurchWrapperJob_RunJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockJob := mock_jobs.NewMockJob(ctrl)
	service := mock_church.NewMockService(ctrl)
	wrapper := newChurchWrapperJob(service, mockJob)
	churchs := buildChurchs()
	t.Run("Success", func(t *testing.T) {
		mockJob.EXPECT().RunJob(gomock.Any()).Return(nil).Times(len(churchs))
		service.EXPECT().List().Return(churchs, nil)
		assert.Nil(t, wrapper.RunJob(context.Background()))
	})
	t.Run("Fail run job", func(t *testing.T) {
		mockJob.EXPECT().RunJob(gomock.Any()).Return(genericError).Times(len(churchs))
		service.EXPECT().List().Return(churchs, nil)
		assert.Nil(t, wrapper.RunJob(context.Background()))
	})
	t.Run("Fail get churchs", func(t *testing.T) {
		service.EXPECT().List().Return(nil, genericError)
		assert.NotNil(t, wrapper.RunJob(context.Background()))
	})
}
