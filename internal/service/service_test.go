package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"awesomeProject/internal/models"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// fakeRepo implements CalendarRepository for testing.
type fakeRepo struct {
	createCalled bool
	updateCalled bool
	deleteCalled bool

	dayCalled   bool
	weekCalled  bool
	monthCalled bool

	closeCalled bool

	lastEvent *models.Event

	// Return controls
	errToReturn    error
	eventsForDay   []models.Event
	eventsForWeek  []models.Event
	eventsForMonth []models.Event
	errForDay      error
	errForWeek     error
	errForMonth    error
	errForCreate   error
	errForUpdate   error
	errForDelete   error
}

func (f *fakeRepo) CreateEvent(ctx context.Context, event *models.Event) error {
	f.createCalled = true
	f.lastEvent = event
	return f.errForCreate
}

func (f *fakeRepo) UpdateEvent(ctx context.Context, event *models.Event) error {
	f.updateCalled = true
	f.lastEvent = event
	return f.errForUpdate
}

func (f *fakeRepo) DeleteEvent(ctx context.Context, event *models.Event) error {
	f.deleteCalled = true
	f.lastEvent = event
	return f.errForDelete
}

func (f *fakeRepo) GetEventsForDay(ctx context.Context, event *models.Event) ([]models.Event, error) {
	f.dayCalled = true
	f.lastEvent = event
	return f.eventsForDay, f.errForDay
}

func (f *fakeRepo) GetEventsForWeek(ctx context.Context, event *models.Event) ([]models.Event, error) {
	f.weekCalled = true
	f.lastEvent = event
	return f.eventsForWeek, f.errForWeek
}

func (f *fakeRepo) GetEventsForMonth(ctx context.Context, event *models.Event) ([]models.Event, error) {
	f.monthCalled = true
	f.lastEvent = event
	return f.eventsForMonth, f.errForMonth
}

func (f *fakeRepo) Close() {
	f.closeCalled = true
}

func newEvent(u int64, name string, d time.Time) *models.Event {
	return &models.Event{
		UserID: u,
		Event:  name,
		Date:   d,
	}
}

// -------- Tests --------

func TestCalendarService_CreateEvent_Success(t *testing.T) {
	r := &fakeRepo{}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	ev := newEvent(42, "Meeting", time.Now())
	err := svc.CreateEvent(context.Background(), ev)
	require.NoError(t, err)
	require.True(t, r.createCalled)
	require.Equal(t, ev, r.lastEvent)
}

func TestCalendarService_CreateEvent_Error(t *testing.T) {
	r := &fakeRepo{errForCreate: errors.New("db failure")}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	ev := newEvent(1, "Fail", time.Now())
	err := svc.CreateEvent(context.Background(), ev)
	require.Error(t, err)
	require.True(t, r.createCalled)
}

func TestCalendarService_UpdateEvent(t *testing.T) {
	r := &fakeRepo{}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	ev := newEvent(7, "UpdateName", time.Now())
	err := svc.UpdateEvent(context.Background(), ev)
	require.NoError(t, err)
	require.True(t, r.updateCalled)
	require.Equal(t, ev, r.lastEvent)
}

func TestCalendarService_UpdateEvent_Error(t *testing.T) {
	r := &fakeRepo{errForUpdate: errors.New("cannot update")}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	ev := newEvent(7, "UpdateName", time.Now())
	err := svc.UpdateEvent(context.Background(), ev)
	require.Error(t, err)
	require.True(t, r.updateCalled)
}

func TestCalendarService_DeleteEvent(t *testing.T) {
	r := &fakeRepo{}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	ev := newEvent(9, "ToDelete", time.Now())
	err := svc.DeleteEvent(context.Background(), ev)
	require.NoError(t, err)
	require.True(t, r.deleteCalled)
	require.Equal(t, ev, r.lastEvent)
}

func TestCalendarService_DeleteEvent_Error(t *testing.T) {
	r := &fakeRepo{errForDelete: errors.New("cannot delete")}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	ev := newEvent(9, "ToDelete", time.Now())
	err := svc.DeleteEvent(context.Background(), ev)
	require.Error(t, err)
	require.True(t, r.deleteCalled)
}

func TestCalendarService_GetEventsForDay(t *testing.T) {
	expected := []models.Event{
		{UserID: 1, Event: "A"},
		{UserID: 1, Event: "B"},
	}
	r := &fakeRepo{eventsForDay: expected}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	ev := newEvent(1, "", time.Now())
	out, err := svc.GetEventsForDay(context.Background(), ev)
	require.NoError(t, err)
	require.True(t, r.dayCalled)
	require.Equal(t, expected, out)
}

func TestCalendarService_GetEventsForDay_Error(t *testing.T) {
	r := &fakeRepo{errForDay: errors.New("query failed")}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	ev := newEvent(1, "", time.Now())
	out, err := svc.GetEventsForDay(context.Background(), ev)
	require.Error(t, err)
	require.Nil(t, out)
	require.True(t, r.dayCalled)
}

func TestCalendarService_GetEventsForWeek(t *testing.T) {
	expected := []models.Event{{UserID: 2, Event: "WeekEvent"}}
	r := &fakeRepo{eventsForWeek: expected}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	ev := newEvent(2, "", time.Now())
	out, err := svc.GetEventsForWeek(context.Background(), ev)
	require.NoError(t, err)
	require.True(t, r.weekCalled)
	require.Equal(t, expected, out)
}

func TestCalendarService_GetEventsForMonth(t *testing.T) {
	expected := []models.Event{{UserID: 3, Event: "MonthEvent1"}, {UserID: 3, Event: "MonthEvent2"}}
	r := &fakeRepo{eventsForMonth: expected}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	ev := newEvent(3, "", time.Now())
	out, err := svc.GetEventsForMonth(context.Background(), ev)
	require.NoError(t, err)
	require.True(t, r.monthCalled)
	require.Equal(t, expected, out)
}

func TestCalendarService_CloseRepo(t *testing.T) {
	r := &fakeRepo{}
	log := zap.NewNop()
	svc := NewCalendarService(r, log)

	svc.CloseRepo()
	require.True(t, r.closeCalled)
}
