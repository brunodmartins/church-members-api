package jobs

import (
	"context"
	"errors"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	mock_member "github.com/brunodmartins/church-members-api/internal/modules/member/mock"
	"github.com/brunodmartins/church-members-api/internal/services/calendar"
	mock_calendar "github.com/brunodmartins/church-members-api/internal/services/calendar/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCalendarJob_truncateYear(t *testing.T) {
	t.Run("Truncate past years", func(t *testing.T) {
		now := time.Now()
		pastDate := now.AddDate(-28, 0, 0)
		assert.Equal(t, now.Format(time.DateOnly), calendarJob{}.truncateYear(pastDate).Format(time.DateOnly))
	})
	t.Run("Truncate current year", func(t *testing.T) {
		now := time.Now()
		assert.Equal(t, now.Format(time.DateOnly), calendarJob{}.truncateYear(now).Format(time.DateOnly))
	})
}

func TestCalendarJob_buildBirthDayEvent(t *testing.T) {
	now := time.Now()
	ctx := buildContext()
	birthDate := now.AddDate(-30, 0, 0)
	expectedEvent := calendar.Event{
		Title:       "🎂 John Doe",
		Description: "Its the birthday of John Doe",
		SenderEmail: "test@example.com",
		SenderName:  "Test Church",
		Time:        birthDate.AddDate(30, 0, 0),
	}
	event := calendarJob{}.buildBirthDayEvent(ctx, &domain.Member{
		Person: &domain.Person{
			BirthDate: birthDate,
			FirstName: "John",
			LastName:  "Doe",
		},
	})
	assert.Equal(t, expectedEvent, event)
}

func TestCalendarJob_buildMarriageDayEvent(t *testing.T) {
	now := time.Now()
	ctx := buildContext()
	marriageDate := now.AddDate(-5, 0, 0)
	expectedEvent := calendar.Event{
		Title:       "🎂 John Doe & Doe John",
		Description: "Its the marriage day of John Doe & Doe John",
		SenderEmail: "test@example.com",
		SenderName:  "Test Church",
		Time:        marriageDate.AddDate(5, 0, 0),
	}
	event := calendarJob{}.buildMarriageDayEvent(ctx, &domain.Member{
		Person: &domain.Person{
			MarriageDate: &marriageDate,
			FirstName:    "John",
			LastName:     "Doe",
			SpousesName:  "Doe John",
		},
	})
	assert.Equal(t, expectedEvent, event)
}

func TestCalendarJob_RunJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := buildContext()

	storage := mock_calendar.NewMockStorage(ctrl)
	membersService := mock_member.NewMockService(ctrl)

	job := NewCalendarJob(storage, membersService)

	t.Run("Run the Job successfully, writing the calendar into the storage", func(t *testing.T) {
		membersService.EXPECT().SearchMembers(gomock.Eq(ctx), gomock.Any()).Return(BuilderMembers(), nil)
		storage.EXPECT().Store(gomock.Eq(ctx), gomock.AssignableToTypeOf(&calendar.Calendar{})).DoAndReturn(func(ctx context.Context, calendar *calendar.Calendar) error {
			assert.Equal(t, "church_short_name_calendar.ics", calendar.Name)
			assert.Equal(t, 2, len(calendar.Events()))
			return nil
		})
		assert.NoError(t, job.RunJob(ctx))
	})
	t.Run("Run the Job with error on writing the calendar into the storage", func(t *testing.T) {
		membersService.EXPECT().SearchMembers(gomock.Eq(ctx), gomock.Any()).Return(BuilderMembers(), nil)
		storage.EXPECT().Store(gomock.Eq(ctx), gomock.AssignableToTypeOf(&calendar.Calendar{})).Return(errors.New("generic error"))
		assert.Error(t, job.RunJob(ctx))
	})
	t.Run("Run the Job with error on searching members", func(t *testing.T) {
		membersService.EXPECT().SearchMembers(gomock.Eq(ctx), gomock.Any()).Return(nil, errors.New("generic error"))
		assert.Error(t, job.RunJob(ctx))
	})
}

func BuilderMembers() []*domain.Member {
	marriageDateJohn2 := time.Now().AddDate(-5, 0, 0)
	return []*domain.Member{
		{
			Person: &domain.Person{
				BirthDate: time.Now(),
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			Person: &domain.Person{
				BirthDate:    time.Now(),
				FirstName:    "John",
				LastName:     "Doe 2",
				MarriageDate: &marriageDateJohn2,
				SpousesName:  "Doe John",
			},
		},
	}
}
