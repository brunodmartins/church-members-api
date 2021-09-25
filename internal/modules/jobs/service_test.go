package jobs

import (
	"fmt"
	mock_email "github.com/BrunoDM2943/church-members-api/internal/services/email/mock"
	"testing"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	mock_member "github.com/BrunoDM2943/church-members-api/internal/modules/member/mock"
	mock_notification "github.com/BrunoDM2943/church-members-api/internal/services/notification/mock"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.Set("bundles.location", "../../../bundles")
	viper.Set("church.shortname", "SHORTNAME")
	viper.Set("jobs.daily.phones", "1,2")
}

func TestLastDayRange(t *testing.T) {
	start, end := lastDaysRange()
	assert.Equal(t, start.Add(6*24*time.Hour), end)
}

func TestFmtDate(t *testing.T) {
	date := time.Date(2020, time.August, 9, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "09-Aug", fmtDate(date))
}

func TestWeeklyBuildMessage(t *testing.T) {
	job := newWeeklyBirthDaysJob(nil, nil)
	time := time.Now()
	fmtDate := fmtDate(time)
	t.Run("With both birth and marriage", func(t *testing.T) {
		expected := fmt.Sprintf("Weekly birthdays\nBirth\n- foo bar - %s\nMarriage\n- foo bar & foo2 bar2 - %s\n", fmtDate, fmtDate)
		assert.Equal(t, expected, job.buildMessage(BuildBirthDaysMembers(time), BuildMarriageMembers(&time)))
	})
	t.Run("Only birth", func(t *testing.T) {
		expected := fmt.Sprintf("Weekly birthdays\nBirth\n- foo bar - %s\nMarriage\n---------\n", fmtDate)
		assert.Equal(t, expected, job.buildMessage(BuildBirthDaysMembers(time), []*domain.Member{}))
	})
	t.Run("Only marriage", func(t *testing.T) {
		expected := fmt.Sprintf("Weekly birthdays\nBirth\n---------\nMarriage\n- foo bar & foo2 bar2 - %s\n", fmtDate)
		assert.Equal(t, expected, job.buildMessage([]*domain.Member{}, BuildMarriageMembers(&time)))
	})
	t.Run("None", func(t *testing.T) {
		expected := "Weekly birthdays\nBirth\n---------\nMarriage\n---------\n"
		assert.Equal(t, expected, job.buildMessage([]*domain.Member{}, []*domain.Member{}))
	})
}

func TestDailyBuildMessage(t *testing.T) {
	job := newDailyBirthDaysJob(nil, nil)
	time := time.Now()
	fmtDate := fmtDate(time)
	expected := fmt.Sprintf("SHORTNAME:Birthdays-foo bar-%s,", fmtDate)
	assert.Equal(t, expected, job.buildMessage(BuildBirthDaysMembers(time)))
}

func TestWeeklyBirthDaysJob_RunJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberService := mock_member.NewMockService(ctrl)
	emailService := mock_email.NewMockService(ctrl)
	job := newWeeklyBirthDaysJob(memberService, emailService)
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		alreadyCalled := false
		memberService.EXPECT().SearchMembers(gomock.Any()).DoAndReturn(func(querySpecification member.QuerySpecification, postSpecification ...member.Specification) ([]*domain.Member, error) {
			if !alreadyCalled {
				alreadyCalled = true
				return BuildBirthDaysMembers(now), nil
			}
			return BuildMarriageMembers(&now), nil
		}).Times(2)
		emailService.EXPECT().SendEmail(gomock.Any()).Return(nil)
		assert.Nil(t, job.RunJob())
	})
	t.Run("Fail Notification", func(t *testing.T) {
		alreadyCalled := false
		memberService.EXPECT().SearchMembers(gomock.Any()).DoAndReturn(func(querySpecification member.QuerySpecification, postSpecification ...member.Specification) ([]*domain.Member, error) {
			if !alreadyCalled {
				alreadyCalled = true
				return BuildBirthDaysMembers(now), nil
			}
			return BuildMarriageMembers(&now), nil
		}).Times(2)
		emailService.EXPECT().SendEmail(gomock.Any()).Return(genericError)
		assert.NotNil(t, job.RunJob())
	})
	t.Run("Fail Search marriage members", func(t *testing.T) {
		alreadyCalled := false
		memberService.EXPECT().SearchMembers(gomock.Any()).DoAndReturn(func(querySpecification member.QuerySpecification, postSpecification ...member.Specification) ([]*domain.Member, error) {
			if !alreadyCalled {
				alreadyCalled = true
				return BuildBirthDaysMembers(now), nil
			}
			return BuildMarriageMembers(&now), genericError
		}).Times(2)
		assert.NotNil(t, job.RunJob())
	})
	t.Run("Fail Search birth members", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any()).Return(nil, genericError)
		assert.NotNil(t, job.RunJob())
	})
}

func TestDailyBirthDaysJob_RunJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	memberService := mock_member.NewMockService(ctrl)
	notificationService := mock_notification.NewMockService(ctrl)
	job := newDailyBirthDaysJob(memberService, notificationService)

	t.Run("Success", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any()).Return(BuildBirthDaysMembers(now), nil)
		notificationService.EXPECT().NotifyMobile(gomock.Any(), gomock.Any()).Return(nil).Times(2)
		assert.Nil(t, job.RunJob())
	})
	t.Run("Fail notify", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any()).Return(BuildBirthDaysMembers(now), nil)
		notificationService.EXPECT().NotifyMobile(gomock.Any(), gomock.Any()).Return(genericError)
		assert.NotNil(t, job.RunJob())
	})
	t.Run("Success - No member", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any()).Return([]*domain.Member{}, nil)
		assert.Nil(t, job.RunJob())
	})
	t.Run("Fail - search members", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any()).Return(nil, genericError)
		assert.NotNil(t, job.RunJob())
	})
}
