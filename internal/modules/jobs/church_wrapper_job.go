package jobs

import (
	"context"
	"github.com/brunodmartins/church-members-api/platform/i18n"
	"golang.org/x/text/language"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/church"
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
	churches, err := wrapper.service.List()
	if err != nil {
		logrus.Error("Error obtaining church list", err)
		return err
	}
	for _, church := range churches {
		err := wrapper.job.RunJob(wrapper.buildContext(ctx, church))
		if err != nil {
			logrus.WithField("church_id", church.ID).Error("Job failed.", err)
		} else {
			logrus.WithField("church_id", church.ID).Info("Job OK.")
		}

	}
	return nil
}

func (wrapper *churchWrapperJob) buildContext(ctx context.Context, church *domain.Church) context.Context {
	ctx = context.WithValue(ctx, "user", &domain.User{
		Church: church,
	})
	return context.WithValue(ctx, "i18n", i18n.GetLocalize(language.MustParse(church.Language)))
}
