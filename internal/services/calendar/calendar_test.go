package calendar

import (
	"bytes"
	"context"
	"errors"
	ics "github.com/arran4/golang-ical"
	mock_storage "github.com/brunodmartins/church-members-api/internal/services/storage/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCalendarBuild(t *testing.T) {
	t.Run("should build calendar", func(t *testing.T) {
		calFile := buildCalendar().Serialize()
		if _, err := ics.ParseCalendar(bytes.NewBuffer(calFile)); err != nil {
			t.Error(err)
		}
	})

}

func buildCalendar() *Calendar {
	testCalendar := New("Test Calendar")
	testCalendar.AddEvent(Event{
		Title:       "Test Event",
		Description: "🎂 " + "Test Event desc",
		SenderEmail: "test@test.com",
		SenderName:  "Test Application",
		Time:        time.Now().Add(-1 * 2 * 24 * time.Hour),
	})
	testCalendar.AddEvent(Event{
		Title:       "Test Event 2",
		Description: "🎂 " + "Test Event desc",
		SenderEmail: "test@test.com",
		SenderName:  "Test Application",
		Time:        time.Now().Add(-1 * 7 * 24 * time.Hour),
	})
	testCalendar.AddEvent(Event{
		Title:       "Test Event 3",
		Description: "🎂 " + "Test Event desc",
		SenderEmail: "test@test.com",
		SenderName:  "Test Application",
		Time:        time.Now().Add(24 * time.Hour),
	})
	return testCalendar
}

func TestCalendarStorage_Store(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()

	t.Run("should store calendar successfully", func(t *testing.T) {
		calendar := buildCalendar()
		mockStorage := mock_storage.NewMockService(ctrl)
		mockStorage.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq(calendar.Name), gomock.AssignableToTypeOf([]byte{})).Return(nil)

		assert.NoError(t, NewCalendarStorage(mockStorage).Store(ctx, calendar))
	})

	t.Run("should fail to store calendar due to storage error", func(t *testing.T) {
		calendar := buildCalendar()
		mockStorage := mock_storage.NewMockService(ctrl)
		mockStorage.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq(calendar.Name), gomock.AssignableToTypeOf([]byte{})).Return(errors.New("generic error"))

		assert.Error(t, NewCalendarStorage(mockStorage).Store(ctx, calendar))
	})
}

func TestCalendarStorage_GetURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()
	name := "test_calendar.ics"

	t.Run("should get calendar URL successfully", func(t *testing.T) {
		expected := "https://fake-url"
		mockStorage := mock_storage.NewMockService(ctrl)
		mockStorage.EXPECT().GetFileURL(gomock.Eq(ctx), gomock.Eq(name)).Return(expected, nil)
		result, err := NewCalendarStorage(mockStorage).GetURL(ctx, name)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should fail to read calendar due to storage error", func(t *testing.T) {
		mockStorage := mock_storage.NewMockService(ctrl)
		mockStorage.EXPECT().GetFileURL(gomock.Eq(ctx), gomock.Eq(name)).Return("", errors.New("generic error"))
		_, err := NewCalendarStorage(mockStorage).GetURL(ctx, name)
		assert.Error(t, err)
	})
}
