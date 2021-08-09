package jobs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildJob(t *testing.T) {
	assert.IsType(t, &dailyBirthDaysJob{}, BuildJob(DAILY_BIRTHDAYS))
	assert.IsType(t, &weeklyBirthDaysJob{}, BuildJob(WEEKLY_BIRTHDAYS))
	assert.Nil(t, BuildJob(10))
}
