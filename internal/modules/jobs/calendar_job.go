package jobs

import (
	"context"
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/services/calendar"
	"github.com/brunodmartins/church-members-api/platform/i18n"
	"time"
)

type calendarJob struct {
	storage       calendar.Storage
	memberService member.Service
}

func (c calendarJob) RunJob(ctx context.Context) error {
	members, err := c.memberService.SearchMembers(ctx, member.OnlyActive())
	if err != nil {
		return err
	}
	churchCalendar := calendar.New(fmt.Sprintf("%s_calendar.ics", domain.GetChurch(ctx).Abbreviation))
	for _, churchMember := range members {
		churchCalendar.AddEvent(c.buildBirthDayEvent(ctx, churchMember))
		if churchMember.Person.IsMarried() {
			churchCalendar.AddEvent(c.buildMarriageDayEvent(ctx, churchMember))
		}
	}
	return c.storage.Store(ctx, churchCalendar)
}

func (c calendarJob) buildBirthDayEvent(ctx context.Context, member *domain.Member) calendar.Event {
	return calendar.Event{
		Title:       fmt.Sprintf("🎂 %s", member.Person.GetFullName()),
		Description: i18n.GetMessage(ctx, "Jobs.Calendar.Birthday.Message", member.Person.GetFullName()),
		SenderEmail: domain.GetChurch(ctx).Email,
		SenderName:  domain.GetChurch(ctx).Name,
		Time:        c.truncateYear(member.Person.BirthDate),
	}
}

func (c calendarJob) buildMarriageDayEvent(ctx context.Context, member *domain.Member) calendar.Event {
	return calendar.Event{
		Title:       fmt.Sprintf("🎂 %s", member.Person.GetCoupleName()),
		Description: i18n.GetMessage(ctx, "Jobs.Calendar.MarriageDay.Message", member.Person.GetCoupleName()),
		SenderEmail: domain.GetChurch(ctx).Email,
		SenderName:  domain.GetChurch(ctx).Name,
		Time:        c.truncateYear(*member.Person.MarriageDate),
	}
}

func (c calendarJob) truncateYear(date time.Time) time.Time {
	currentYear := time.Now().Year()
	dateYear := date.Year()
	return date.AddDate(currentYear-dateYear, 0, 0)
}

func NewCalendarJob(storage calendar.Storage, memberService member.Service) Job {
	return &calendarJob{
		storage:       storage,
		memberService: memberService,
	}
}
