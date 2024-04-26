package calendar

import (
	"context"
	"encoding/base64"
	"fmt"
	ics "github.com/arran4/golang-ical"
	"github.com/brunodmartins/church-members-api/internal/services/storage"
	"time"
)

type Calendar struct {
	Name   string
	events []Event
}

func New(name string) *Calendar {
	return &Calendar{name, make([]Event, 0)}
}

func (calendar *Calendar) AddEvent(event Event) {
	calendar.events = append(calendar.events, event)
}

func (calendar *Calendar) Events() []Event {
	return calendar.events
}

func (calendar *Calendar) Serialize() []byte {
	cal := ics.NewCalendar()
	cal.SetName(calendar.Name)
	cal.SetMethod(ics.MethodRequest)
	for _, event := range calendar.events {
		addEvent := cal.AddEvent(base64.StdEncoding.EncodeToString([]byte(event.Title)))
		addEvent.SetCreatedTime(time.Now())
		addEvent.SetDtStampTime(time.Now())
		addEvent.SetAllDayStartAt(event.Time)
		addEvent.SetAllDayEndAt(event.Time)
		addEvent.SetSummary(event.Title)
		addEvent.SetDescription(event.Description)
		addEvent.AddRrule(fmt.Sprintf("FREQ=YEARLY;BYMONTH=%d;BYMONTHDAY=%d", time.Now().Month(), time.Now().Day()))
		alarm := addEvent.AddAlarm()
		alarm.SetAction(ics.ActionDisplay)
		alarm.SetDescription(event.Description)
		alarm.SetTrigger("P0DT9H0M0S")
		addEvent.SetOrganizer(event.SenderEmail, ics.WithCN(event.SenderName))
	}
	return []byte(cal.Serialize())
}

type Event struct {
	Title       string
	Description string
	SenderEmail string
	SenderName  string
	Time        time.Time
}

//go:generate mockgen -source=./calendar.go -destination=./mock/calendar_mock.go
type Storage interface {
	Store(ctx context.Context, calendar *Calendar) error
	GetURL(ctx context.Context, name string) (string, error)
}

type calendarStorage struct {
	storage storage.Service
}

func (cal *calendarStorage) Store(ctx context.Context, calendar *Calendar) error {
	return cal.storage.SaveFile(ctx, calendar.Name, calendar.Serialize())
}

func (cal *calendarStorage) GetURL(ctx context.Context, name string) (string, error) {
	return cal.storage.GetFileURL(ctx, name)
}

func NewCalendarStorage(storage storage.Service) Storage {
	return &calendarStorage{
		storage: storage,
	}
}
