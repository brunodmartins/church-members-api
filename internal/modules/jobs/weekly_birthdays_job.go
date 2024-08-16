package jobs

import (
	"context"
	_ "embed"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/brunodmartins/church-members-api/internal/services/email"
	"github.com/brunodmartins/church-members-api/platform/i18n"
	"sort"
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
	emailData := job.buildEmailData(ctx, birthMembers, marriageMembers)
	users, err := job.userService.SearchUser(ctx, user.WithEmailNotifications())
	if err != nil {
		return err
	}
	for _, emailTO := range mapToSlice(users, getEmail) {
		if err := job.emailService.SendTemplateEmail(email.WeeklyBirthTemplate, emailData, i18n.GetMessage(ctx, "Emails.WeeklyBirth.Title"), emailTO); err != nil {
			return err
		}
	}
	return nil
}

func (job weeklyBirthDaysJob) buildEmailData(ctx context.Context, birthMembers, marriageMembers []*domain.Member) email.WeeklyBirthTemplateDTO {
	weeklyTemplate := email.NewWeeklyBirthTemplateDTO(ctx)
	for _, member := range birthMembers {
		weeklyTemplate.MembersBirth = append(weeklyTemplate.MembersBirth, email.MemberDTO{
			Name: member.Person.GetFullName(),
			Date: fmtDate(member.Person.BirthDate),
		})
	}

	for _, member := range marriageMembers {
		weeklyTemplate.MembersMarriage = append(weeklyTemplate.MembersMarriage, email.MemberDTO{
			Name: member.Person.GetCoupleName(),
			Date: fmtDate(*member.Person.MarriageDate),
		})
	}
	return weeklyTemplate
}

func (weeklyBirthDaysJob) lastDaysRange() (time.Time, time.Time) {
	now := time.Now()
	return now.Add(-6 * 24 * time.Hour), now
}
