package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildJob(t *testing.T) {
	assert.IsType(t, &dailyBirthDaysJob{}, BuildJob(DAILY_BIRTHDAYS).(*churchWrapperJob).job)
	assert.IsType(t, &weeklyBirthDaysJob{}, BuildJob(WEEKLY_BIRTHDAYS).(*churchWrapperJob).job)
	assert.IsType(t, &calendarJob{}, BuildJob(CALENDAR_BUILD).(*churchWrapperJob).job)
	assert.Nil(t, BuildJob(10))
}
