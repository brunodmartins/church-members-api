package jobs

import (
	"fmt"
	"strings"
)

type JobType int

const (
	WEEKLY_BIRTHDAYS JobType = iota + 1
	DAILY_BIRTHDAYS
	CALENDAR_BUILD
)

var jobs = []string{"WEEKLY_BIRTHDAYS", "DAILY_BIRTHDAYS", "CALENDAR_BUILD"}

func (j JobType) String() string {
	return jobs[j-1]
}

func (JobType) From(value string) (JobType, error) {
	value = strings.ToUpper(value)
	for i := 0; i < len(jobs); i++ {
		if strings.ToUpper(jobs[i]) == value {
			return JobType(i + 1), nil
		}
	}
	return 0, fmt.Errorf("invalid job: %s", value)
}
