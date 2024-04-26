package jobs

import (
	"context"
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/brunodmartins/church-members-api/internal/services/email"
	"github.com/brunodmartins/church-members-api/platform/i18n"
	"sort"
	"strings"
	"time"
)

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
	birthMembers, err := job.memberService.SearchMembers(ctx, member.LastBirths(job.lastDaysRange()))
	if err != nil {
		return err
	}
	sort.Sort(domain.SortByBirthDay(birthMembers))
	marriageMembers, err := job.memberService.SearchMembers(ctx, member.LastMarriages(job.lastDaysRange()))
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
		if err := job.emailService.SendEmail(job.buildEmailCommand(ctx, emailBody, emailTO)); err != nil {
			return err
		}
	}
	return nil
}

func (job weeklyBirthDaysJob) buildEmailCommand(ctx context.Context, message, emailTO string) email.Command {
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

func (weeklyBirthDaysJob) lastDaysRange() (time.Time, time.Time) {
	now := time.Now()
	return now.Add(-6 * 24 * time.Hour), now
}
