package jobs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJobType_From(t *testing.T) {
	t.Run("DAILY_BIRTHDAYS", func(t *testing.T) {
		jobType, _ := new(JobType).From("DAILY_BIRTHDAYS")
		assert.Equal(t, DAILY_BIRTHDAYS, jobType)
	})
	t.Run("WEEKLY_BIRTHDAYS", func(t *testing.T) {
		jobType, _ := new(JobType).From("WEEKLY_BIRTHDAYS")
		assert.Equal(t, WEEKLY_BIRTHDAYS, jobType)
	})
	t.Run("CALENDAR_BUILD", func(t *testing.T) {
		jobType, _ := new(JobType).From("CALENDAR_BUILD")
		assert.Equal(t, CALENDAR_BUILD, jobType)
	})
	t.Run("error", func(t *testing.T) {
		_, err := new(JobType).From("err")
		assert.NotNil(t, err)
	})
}

func TestJobType_String(t *testing.T) {
	assert.Equal(t, "DAILY_BIRTHDAYS", DAILY_BIRTHDAYS.String())
	assert.Equal(t, "WEEKLY_BIRTHDAYS", WEEKLY_BIRTHDAYS.String())
}
