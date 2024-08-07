package jobs

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/brunodmartins/church-members-api/internal/services/email"
	"github.com/brunodmartins/church-members-api/platform/i18n"
	"html/template"
	"sort"
	"time"
)

//go:embed resources/weekly_birthdays_template.html
var emailHTML string

type weeklyBirthDaysJob struct {
	memberService member.Service
	emailService  email.Service
	userService   user.Service
}

type weeklyTemplateDTO struct {
	Title           string
	BirthTitle      string
	MarriageTitle   string
	NameColumn      string
	DateColumn      string
	MembersBirth    []memberDTO
	MembersMarriage []memberDTO
}

type memberDTO struct {
	Name string
	Date string
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
	emailBody, err := job.buildMessage(ctx, birthMembers, marriageMembers)
	if err != nil {
		return err
	}
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

func (job weeklyBirthDaysJob) buildMessage(ctx context.Context, birthMembers, marriageMembers []*domain.Member) (string, error) {
	weeklyTemplate := weeklyTemplateDTO{
		Title:         i18n.GetMessage(ctx, "Jobs.Weekly.Title"),
		BirthTitle:    i18n.GetMessage(ctx, "Jobs.Weekly.Birth"),
		MarriageTitle: i18n.GetMessage(ctx, "Jobs.Weekly.Marriage"),
		NameColumn:    i18n.GetMessage(ctx, "Domain.Name"),
		DateColumn:    i18n.GetMessage(ctx, "Domain.Date"),
	}
	tmpl, err := template.New("weekly_job").Parse(emailHTML)
	if err != nil {
		return "", err
	}
	for _, member := range birthMembers {
		weeklyTemplate.MembersBirth = append(weeklyTemplate.MembersBirth, memberDTO{
			Name: member.Person.GetFullName(),
			Date: fmtDate(member.Person.BirthDate),
		})
	}

	for _, member := range marriageMembers {
		weeklyTemplate.MembersMarriage = append(weeklyTemplate.MembersMarriage, memberDTO{
			Name: member.Person.GetCoupleName(),
			Date: fmtDate(*member.Person.MarriageDate),
		})
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, weeklyTemplate)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (weeklyBirthDaysJob) lastDaysRange() (time.Time, time.Time) {
	now := time.Now()
	return now.Add(-6 * 24 * time.Hour), now
}
