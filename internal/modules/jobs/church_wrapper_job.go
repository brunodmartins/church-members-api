package jobs

import (
	"context"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/modules/church"
	"github.com/sirupsen/logrus"
)

type churchWrapperJob struct {
	service church.Service
	job     Job
}

func newChurchWrapperJob(service church.Service, job Job) Job {
	return &churchWrapperJob{
		service: service,
		job:     job,
	}
}

func (wrapper *churchWrapperJob) RunJob(ctx context.Context) error {
	churchs, err := wrapper.service.List()
	if err != nil {
		logrus.Error("Error obtaining church list", err)
		return err
	}
	for _, church := range churchs {
		err := wrapper.job.RunJob(context.WithValue(ctx, "user", &domain.User{
			ChurchID: church.ID,
		}))
		if err != nil {
			logrus.WithField("church_id", church.ID).Error("Job failed.", err)
		} else {
			logrus.WithField("church_id", church.ID).Info("Job OK.")
		}

	}
	return nil
}
