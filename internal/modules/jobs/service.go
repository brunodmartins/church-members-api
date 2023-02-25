package jobs

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/brunodmartins/church-members-api/internal/services/email"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/services/notification"
	"github.com/brunodmartins/church-members-api/platform/i18n"
)

// Job exposing jobs operations
//
//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Job interface {
	RunJob(ctx context.Context) error
}

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

type weeklyBirthDaysJob struct {
	memberService member.Service
	emailService  email.Service
	userService   user.Service
}

func newWeeklyBirthDaysJob(
	memberService member.Service,
	emailService email.Service,
	userService user.Service) *weeklyBirthDaysJob {
	return &weeklyBirthDaysJob{
		memberService: memberService,
		emailService:  emailService,
		userService:   userService}
}

func (job weeklyBirthDaysJob) RunJob(ctx context.Context) error {
	birthMembers, err := job.memberService.SearchMembers(ctx, member.LastBirths(lastDaysRange()))
	if err != nil {
		return err
	}
	sort.Sort(domain.SortByBirthDay(birthMembers))
	marriageMembers, err := job.memberService.SearchMembers(ctx, member.LastMarriages(lastDaysRange()))
	if err != nil {
		return err
	}
	sort.Sort(domain.SortByMarriageDay(marriageMembers))
	emailBody := job.buildMessage(ctx, birthMembers, marriageMembers)
	users, err := job.userService.SearchUser(ctx, user.WithEmailNotifications())
	if err != nil {
		return err
	}
	for _, emailTO := range mapToSlice(users, getEmail) {
		if err := job.emailService.SendEmail(buildEmailCommand(ctx, emailBody, emailTO)); err != nil {
			return err
		}
	}
	return nil
}

func buildEmailCommand(ctx context.Context, message, emailTO string) email.Command {
	return email.Command{
		Body:       message,
		Subject:    i18n.GetMessage(ctx, "Jobs.Weekly.Title"),
		Recipients: []string{emailTO},
	}
}

func (job weeklyBirthDaysJob) buildMessage(ctx context.Context, birthMembers, marriageMembers []*domain.Member) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", i18n.GetMessage(ctx, "Jobs.Weekly.Title")))
	builder.WriteString(fmt.Sprintf("%s\n", i18n.GetMessage(ctx, "Jobs.Weekly.Birth")))
	for _, member := range birthMembers {
		builder.WriteString(fmt.Sprintf("- %s - %s\n", member.Person.GetFullName(), fmtDate(member.Person.BirthDate)))
	}
	if len(birthMembers) == 0 {
		builder.WriteString("---------\n")
	}

	builder.WriteString(fmt.Sprintf("%s\n", i18n.GetMessage(ctx, "Jobs.Weekly.Marriage")))
	for _, member := range marriageMembers {
		builder.WriteString(fmt.Sprintf("- %s & %s - %s\n", member.Person.GetFullName(), member.Person.SpousesName, fmtDate(*member.Person.MarriageDate)))
	}
	if len(marriageMembers) == 0 {
		builder.WriteString("---------\n")
	}
	return builder.String()
}

func fmtDate(date time.Time) string {
	return date.Format("02-Jan")
}

func lastDaysRange() (time.Time, time.Time) {
	now := time.Now()
	return now.Add(-6 * 24 * time.Hour), now
}

func getEmail(user *domain.User) string {
	return user.Email
}

func getPhone(user *domain.User) string {
	return user.Phone
}

func mapToSlice(users []*domain.User, mapper func(user *domain.User) string) []string {
	var result []string
	for _, user := range users {
		result = append(result, mapper(user))
	}
	return result
}
