package jobs

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	mock_member "github.com/BrunoDM2943/church-members-api/internal/modules/member/mock"
	mock_notification "github.com/BrunoDM2943/church-members-api/internal/services/notification/mock"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	viper.Set("bundles.location", "../../../bundles")
	viper.Set("church.shortname", "SHORTNAME")
	viper.Set("jobs.daily.phones", "1,2")
}

func TestLastDayRange(t *testing.T) {
	start, end := lastDaysRange()
	assert.Equal(t, start.Add(7*24*time.Hour), end)
}

func TestFmtDate(t *testing.T) {
	date := time.Date(2020, time.August, 9, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "09-Aug", fmtDate(date))
}

func TestWeeklyBuildMessage(t *testing.T) {
	job := newWeeklyBirthDaysJob(nil, nil)
	t.Run("With both birth and marriage", func(t *testing.T) {
		expected := "Weekly birthdays\nBirth\n- foo bar - 09-Aug\nMarriage\n- foo bar & foo2 bar2 - 09-Aug\n"
		assert.Equal(t, expected, job.buildMessage(BuildBirthDaysMembers(), BuildMarriageMembers()))
	})
	t.Run("Only birth", func(t *testing.T) {
		expected := "Weekly birthdays\nBirth\n- foo bar - 09-Aug\nMarriage\n---------\n"
		assert.Equal(t, expected, job.buildMessage(BuildBirthDaysMembers(), []*domain.Member{}))
	})
	t.Run("Only marriage", func(t *testing.T) {
		expected := "Weekly birthdays\nBirth\n---------\nMarriage\n- foo bar & foo2 bar2 - 09-Aug\n"
		assert.Equal(t, expected, job.buildMessage([]*domain.Member{}, BuildMarriageMembers()))
	})
	t.Run("None", func(t *testing.T) {
		expected := "Weekly birthdays\nBirth\n---------\nMarriage\n---------\n"
		assert.Equal(t, expected, job.buildMessage([]*domain.Member{}, []*domain.Member{}))
	})
}

func TestDailyBuildMessage(t *testing.T) {
	job := newDailyBirthDaysJob(nil, nil)
	expected := "SHORTNAME:Birthdays-foo bar-09-Aug,"
	assert.Equal(t, expected, job.buildMessage(BuildBirthDaysMembers()))
}

func TestWeeklyBirthDaysJob_RunJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberService := mock_member.NewMockService(ctrl)
	notificationService := mock_notification.NewMockService(ctrl)
	job := newWeeklyBirthDaysJob(memberService, notificationService)

	t.Run("Success", func(t *testing.T) {
		alreadyCalled := false
		memberService.EXPECT().SearchMembers(gomock.Any()).DoAndReturn(func(querySpecification member.QuerySpecification, postSpecification ...member.Specification) ([]*domain.Member, error) {
			if !alreadyCalled {
				alreadyCalled = true
				return BuildBirthDaysMembers(), nil
			}
			return BuildMarriageMembers(), nil
		}).Times(2)
		notificationService.EXPECT().NotifyTopic(gomock.Any()).Return(nil)
		assert.Nil(t, job.RunJob())
	})
	t.Run("Fail Notification", func(t *testing.T) {
		alreadyCalled := false
		memberService.EXPECT().SearchMembers(gomock.Any()).DoAndReturn(func(querySpecification member.QuerySpecification, postSpecification ...member.Specification) ([]*domain.Member, error) {
			if !alreadyCalled {
				alreadyCalled = true
				return BuildBirthDaysMembers(), nil
			}
			return BuildMarriageMembers(), nil
		}).Times(2)
		notificationService.EXPECT().NotifyTopic(gomock.Any()).Return(genericError)
		assert.NotNil(t, job.RunJob())
	})
	t.Run("Fail Search marriage members", func(t *testing.T) {
		alreadyCalled := false
		memberService.EXPECT().SearchMembers(gomock.Any()).DoAndReturn(func(querySpecification member.QuerySpecification, postSpecification ...member.Specification) ([]*domain.Member, error) {
			if !alreadyCalled {
				alreadyCalled = true
				return BuildBirthDaysMembers(), nil
			}
			return BuildMarriageMembers(), genericError
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

	memberService := mock_member.NewMockService(ctrl)
	notificationService := mock_notification.NewMockService(ctrl)
	job := newDailyBirthDaysJob(memberService, notificationService)

	t.Run("Success", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any()).Return(BuildBirthDaysMembers(), nil)
		notificationService.EXPECT().NotifyMobile(gomock.Any(), gomock.Any()).Return(nil).Times(2)
		assert.Nil(t, job.RunJob())
	})
	t.Run("Fail notify", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any()).Return(BuildBirthDaysMembers(), nil)
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
