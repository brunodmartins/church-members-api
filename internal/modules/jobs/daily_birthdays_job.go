package jobs

import (
	"context"
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/brunodmartins/church-members-api/internal/services/notification"
	"github.com/brunodmartins/church-members-api/platform/i18n"
	"strings"
	"time"
)

type dailyBirthDaysJob struct {
	memberService       member.Service
	userService         user.Service
	notificationService notification.Service
}

func newDailyBirthDaysJob(
	memberService member.Service,
	notificationService notification.Service,
	userService user.Service) *dailyBirthDaysJob {
	return &dailyBirthDaysJob{
		memberService:       memberService,
		notificationService: notificationService,
		userService:         userService,
	}
}

func (job dailyBirthDaysJob) RunJob(ctx context.Context) error {
	members, err := job.memberService.SearchMembers(ctx, member.WithBirthday(time.Now()))
	if err != nil {
		return err
	}
	if len(members) == 0 {
		return nil
	}
	message := job.buildMessage(ctx, members)
	users, err := job.userService.SearchUser(ctx, user.WithSMSNotifications())
	if err != nil {
		return err
	}
	for _, phone := range mapToSlice(users, getPhone) {
		if err := job.notificationService.NotifyMobile(message, phone); err != nil {
			return err
		}
	}
	return nil
}

func (job dailyBirthDaysJob) buildMessage(ctx context.Context, members []*domain.Member) string {
	title := i18n.GetMessage(ctx, "Jobs.Daily.Title")
	builder := strings.Builder{}
	for _, member := range members {
		builder.WriteString(fmt.Sprintf("%s-%s,", member.Person.GetFullName(), fmtDate(member.Person.BirthDate)))
	}
	return fmt.Sprintf("%s:%s-%s", domain.GetChurch(ctx).Abbreviation, title, builder.String())
}
