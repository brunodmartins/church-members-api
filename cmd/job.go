package cmd

import (
	"context"
	"errors"

	"github.com/BrunoDM2943/church-members-api/internal/modules/jobs"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

type JobApplication struct{}

type JobEvent struct {
	Name string `json:"name"`
}

func (JobApplication) Run() {
	lambda.Start(func(ctx context.Context, event JobEvent) error {
		logrus.Infof("Received event: %v", event)
		jobType, err := new(jobs.JobType).From(event.Name)
		if err != nil {
			return err
		}
		job := jobs.BuildJob(jobType)
		if job == nil {
			return errors.New("job not found")
		}
		return job.RunJob()
	})
}
