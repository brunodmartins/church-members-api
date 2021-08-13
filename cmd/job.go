package cmd

import (
	"encoding/json"
	"errors"

	"github.com/BrunoDM2943/church-members-api/internal/modules/jobs"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

type JobApplication struct{}

func (JobApplication) Run() {
	lambda.Start(func(payload []byte) error {
		var input map[string]interface{}
		logrus.Infof("Received event: %s", string(payload))
		err := json.Unmarshal(payload, &input)
		if err != nil {
			return err
		}
		jobType, err := new(jobs.JobType).From(input["name"].(string))
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
