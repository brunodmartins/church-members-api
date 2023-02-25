package jobs

import (
	"context"
	"fmt"
	"testing"
	"time"

	mock_user "github.com/brunodmartins/church-members-api/internal/modules/user/mock"
	mock_email "github.com/brunodmartins/church-members-api/internal/services/email/mock"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	mock_member "github.com/brunodmartins/church-members-api/internal/modules/member/mock"
	mock_notification "github.com/brunodmartins/church-members-api/internal/services/notification/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLastDayRange(t *testing.T) {
	start, end := lastDaysRange()
	assert.Equal(t, start.Add(6*24*time.Hour), end)
}

func TestFmtDate(t *testing.T) {
	date := time.Date(2020, time.August, 9, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "09-Aug", fmtDate(date))
}

func TestWeeklyBuildMessage(t *testing.T) {
	job := newWeeklyBirthDaysJob(nil, nil, nil)
	now := time.Now()
	fmtDate := fmtDate(now)
	t.Run("With both birth and marriage", func(t *testing.T) {
		expected := fmt.Sprintf("Weekly birthdays\nBirth\n- foo bar - %s\nMarriage\n- foo bar & foo2 bar2 - %s\n", fmtDate, fmtDate)
		assert.Equal(t, expected, job.buildMessage(context.TODO(), BuildBirthDaysMembers(now), BuildMarriageMembers(&now)))
	})
	t.Run("Only birth", func(t *testing.T) {
		expected := fmt.Sprintf("Weekly birthdays\nBirth\n- foo bar - %s\nMarriage\n---------\n", fmtDate)
		assert.Equal(t, expected, job.buildMessage(context.TODO(), BuildBirthDaysMembers(now), []*domain.Member{}))
	})
	t.Run("Only marriage", func(t *testing.T) {
		expected := fmt.Sprintf("Weekly birthdays\nBirth\n---------\nMarriage\n- foo bar & foo2 bar2 - %s\n", fmtDate)
		assert.Equal(t, expected, job.buildMessage(context.TODO(), []*domain.Member{}, BuildMarriageMembers(&now)))
	})
	t.Run("None", func(t *testing.T) {
		expected := "Weekly birthdays\nBirth\n---------\nMarriage\n---------\n"
		assert.Equal(t, expected, job.buildMessage(context.TODO(), []*domain.Member{}, []*domain.Member{}))
	})
}

func TestDailyBuildMessage(t *testing.T) {
	job := newDailyBirthDaysJob(nil, nil, nil)
	now := time.Now()
	fmtDate := fmtDate(now)
	expected := fmt.Sprintf("church_short_name:Birthdays-foo bar-%s,", fmtDate)
	assert.Equal(t, expected, job.buildMessage(buildContext(), BuildBirthDaysMembers(now)))
}

func TestWeeklyBirthDaysJob_RunJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberService := mock_member.NewMockService(ctrl)
	emailService := mock_email.NewMockService(ctrl)
	userService := mock_user.NewMockService(ctrl)
	job := newWeeklyBirthDaysJob(memberService, emailService, userService)
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		alreadyCalled := false
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, querySpecification wrapper.QuerySpecification, postSpecification ...member.Specification) ([]*domain.Member, error) {
			if !alreadyCalled {
				alreadyCalled = true
				return BuildBirthDaysMembers(now), nil
			}
			return BuildMarriageMembers(&now), nil
		}).Times(2)
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(BuildUsers(), nil)
		emailService.EXPECT().SendEmail(gomock.Any()).Return(nil).Times(2)
		assert.Nil(t, job.RunJob(buildContext()))
	})
	t.Run("Fail users search", func(t *testing.T) {
		alreadyCalled := false
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, querySpecification wrapper.QuerySpecification, postSpecification ...member.Specification) ([]*domain.Member, error) {
			if !alreadyCalled {
				alreadyCalled = true
				return BuildBirthDaysMembers(now), nil
			}
			return BuildMarriageMembers(&now), nil
		}).Times(2)
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(BuildUsers(), genericError)
		assert.NotNil(t, job.RunJob(buildContext()))
	})
	t.Run("Fail Notification", func(t *testing.T) {
		alreadyCalled := false
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, querySpecification wrapper.QuerySpecification, postSpecification ...member.Specification) ([]*domain.Member, error) {
			if !alreadyCalled {
				alreadyCalled = true
				return BuildBirthDaysMembers(now), nil
			}
			return BuildMarriageMembers(&now), nil
		}).Times(2)
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(BuildUsers(), nil)
		emailService.EXPECT().SendEmail(gomock.Any()).Return(genericError)
		assert.NotNil(t, job.RunJob(buildContext()))
	})
	t.Run("Fail Search marriage members", func(t *testing.T) {
		alreadyCalled := false
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, querySpecification wrapper.QuerySpecification, postSpecification ...member.Specification) ([]*domain.Member, error) {
			if !alreadyCalled {
				alreadyCalled = true
				return BuildBirthDaysMembers(now), nil
			}
			return BuildMarriageMembers(&now), genericError
		}).Times(2)
		assert.NotNil(t, job.RunJob(buildContext()))
	})
	t.Run("Fail Search birth members", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).Return(nil, genericError)
		assert.NotNil(t, job.RunJob(buildContext()))
	})
}

func TestDailyBirthDaysJob_RunJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	memberService := mock_member.NewMockService(ctrl)
	notificationService := mock_notification.NewMockService(ctrl)
	userService := mock_user.NewMockService(ctrl)
	job := newDailyBirthDaysJob(memberService, notificationService, userService)

	t.Run("Success", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).Return(BuildBirthDaysMembers(now), nil)
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(BuildUsers(), nil)
		notificationService.EXPECT().NotifyMobile(gomock.Any(), gomock.Any()).Return(nil).Times(2)
		assert.Nil(t, job.RunJob(buildContext()))
	})
	t.Run("Fail users search", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).Return(BuildBirthDaysMembers(now), nil)
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(BuildUsers(), genericError)
		assert.NotNil(t, job.RunJob(buildContext()))
	})
	t.Run("Fail notify", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).Return(BuildBirthDaysMembers(now), nil)
		userService.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(BuildUsers(), nil)
		notificationService.EXPECT().NotifyMobile(gomock.Any(), gomock.Any()).Return(genericError)
		assert.NotNil(t, job.RunJob(buildContext()))
	})
	t.Run("Success - No member", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).Return([]*domain.Member{}, nil)
		assert.Nil(t, job.RunJob(buildContext()))
	})
	t.Run("Fail - search members", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.Any()).Return(nil, genericError)
		assert.NotNil(t, job.RunJob(buildContext()))
	})
}

func buildContext() context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{
			ID:           "church_id_test",
			Abbreviation: "church_short_name",
		},
	})
}
