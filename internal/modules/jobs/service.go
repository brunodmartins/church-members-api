package jobs

import (
	"fmt"
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
	notificationService notification.Service
}

func newDailyBirthDaysJob(memberService member.Service, notificationService notification.Service) *dailyBirthDaysJob {
	return &dailyBirthDaysJob{memberService: memberService, notificationService: notificationService}
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
	for _, phone := range strings.Split(viper.GetString("jobs.daily.phones"), ",") {
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
	notificationService notification.Service
}

func newWeeklyBirthDaysJob(memberService member.Service, notificationService notification.Service) *weeklyBirthDaysJob {
	return &weeklyBirthDaysJob{memberService: memberService, notificationService: notificationService}
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
	return job.notificationService.NotifyTopic(job.buildMessage(birthMembers, marriageMembers))
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
