package jobs

import (
	"fmt"
	"github.com/BrunoDM2943/church-members-api/internal/modules/user"
	"github.com/BrunoDM2943/church-members-api/internal/services/email"
	"sort"
	"strings"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/BrunoDM2943/church-members-api/internal/services/notification"
	"github.com/BrunoDM2943/church-members-api/platform/i18n"
	"github.com/spf13/viper"
)

//Job exposing jobs operations
type Job interface {
	RunJob() error
}

type dailyBirthDaysJob struct {
	memberService       member.Service
	userService			user.Service
	notificationService notification.Service
}

func newDailyBirthDaysJob(
	memberService member.Service,
	notificationService notification.Service,
	userService user.Service) *dailyBirthDaysJob {
	return &dailyBirthDaysJob{
		memberService: memberService,
		notificationService: notificationService,
		userService: userService,
	}
}

func (job dailyBirthDaysJob) RunJob() error {
	members, err := job.memberService.SearchMembers(member.WithBirthday(time.Now()))
	if err != nil {
		return err
	}
	if len(members) == 0 {
		return nil
	}
	message := job.buildMessage(members)
	users, err := job.userService.SearchUser(user.WithSMSNotifications())
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

func (job dailyBirthDaysJob) buildMessage(members []*domain.Member) string {
	tr := i18n.GetMessageService()
	title := tr.GetMessage("Jobs.Daily.Title", "Birthdays")
	churchName := viper.GetString("church.shortname")
	builder := strings.Builder{}
	for _, member := range members {
		builder.WriteString(fmt.Sprintf("%s-%s,", member.Person.GetFullName(), fmtDate(member.Person.BirthDate)))
	}
	return fmt.Sprintf("%s:%s-%s", churchName, title, builder.String())
}

type weeklyBirthDaysJob struct {
	memberService       member.Service
	emailService 		email.Service
	userService 		user.Service
}

func newWeeklyBirthDaysJob(
	memberService member.Service,
	emailService email.Service,
	userService user.Service) *weeklyBirthDaysJob {
	return &weeklyBirthDaysJob{
		memberService: memberService,
		emailService: emailService,
		userService: userService}
}

func (job weeklyBirthDaysJob) RunJob() error {
	birthMembers, err := job.memberService.SearchMembers(member.LastBirths(lastDaysRange()))
	if err != nil {
		return err
	}
	sort.Sort(domain.SortByBirthDay(birthMembers))
	marriageMembers, err := job.memberService.SearchMembers(member.LastMarriages(lastDaysRange()))
	if err != nil {
		return err
	}
	sort.Sort(domain.SortByMarriageDay(marriageMembers))
	emailBody := job.buildMessage(birthMembers, marriageMembers)
	users, err := job.userService.SearchUser(user.WithEmailNotifications())
	if err != nil {
		return err
	}
	for _, emailTO := range mapToSlice(users, getEmail) {
		if err := job.emailService.SendEmail(buildEmailCommand(emailBody, emailTO)); err != nil {
			return err
		}
	}
	return nil
}

func buildEmailCommand(message, emailTO string) email.Command {
	tr := i18n.GetMessageService()
	return email.Command{
		Body: message,
		Subject: tr.GetMessage("Jobs.Weekly.Title", "Weekly birthdays"),
		Recipients: []string{emailTO},
	}
}

func (job weeklyBirthDaysJob) buildMessage(birthMembers, marriageMembers []*domain.Member) string {
	tr := i18n.GetMessageService()
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", tr.GetMessage("Jobs.Weekly.Title", "Weekly birthdays")))
	builder.WriteString(fmt.Sprintf("%s\n", tr.GetMessage("Jobs.Weekly.Birth", "Birth")))
	for _, member := range birthMembers {
		builder.WriteString(fmt.Sprintf("- %s - %s\n", member.Person.GetFullName(), fmtDate(member.Person.BirthDate)))
	}
	if len(birthMembers) == 0 {
		builder.WriteString("---------\n")
	}

	builder.WriteString(fmt.Sprintf("%s\n", tr.GetMessage("Jobs.Weekly.Marriage", "Marriage")))
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