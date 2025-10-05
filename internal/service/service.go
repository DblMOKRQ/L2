package service

import (
	"awesomeProject/internal/models"
	"context"
	"go.uber.org/zap"
)

type CalendarRepository interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, event *models.Event) error
	GetEventsForDay(ctx context.Context, event *models.Event) ([]models.Event, error)
	GetEventsForWeek(ctx context.Context, event *models.Event) ([]models.Event, error)
	GetEventsForMonth(ctx context.Context, event *models.Event) ([]models.Event, error)
	Close()
}

type CalendarService struct {
	repo CalendarRepository
	log  *zap.Logger
}

func NewCalendarService(repo CalendarRepository, log *zap.Logger) *CalendarService {
	return &CalendarService{repo: repo, log: log.Named("CalendarService")}
}

func (s *CalendarService) CreateEvent(ctx context.Context, event *models.Event) error {
	s.log.Info("Creating event", zap.Int64("user_id", event.UserID), zap.Time("date", event.Date), zap.String("event", event.Event))
	return s.repo.CreateEvent(ctx, event)
}
func (s *CalendarService) UpdateEvent(ctx context.Context, event *models.Event) error {
	s.log.Info("Updating event", zap.Int64("user_id", event.UserID), zap.Time("date", event.Date), zap.String("event", event.Event))
	return s.repo.UpdateEvent(ctx, event)
}

func (s *CalendarService) DeleteEvent(ctx context.Context, event *models.Event) error {
	s.log.Info("Deleting event", zap.Int64("user_id", event.UserID), zap.String("event", event.Event))
	return s.repo.DeleteEvent(ctx, event)
}

func (s *CalendarService) GetEventsForDay(ctx context.Context, event *models.Event) ([]models.Event, error) {
	s.log.Info("Getting events for day", zap.Int64("user_id", event.UserID), zap.Time("date", event.Date))
	return s.repo.GetEventsForDay(ctx, event)
}
func (s *CalendarService) GetEventsForWeek(ctx context.Context, event *models.Event) ([]models.Event, error) {
	s.log.Info("Getting events for week", zap.Int64("user_id", event.UserID))
	return s.repo.GetEventsForWeek(ctx, event)
}
func (s *CalendarService) GetEventsForMonth(ctx context.Context, event *models.Event) ([]models.Event, error) {
	s.log.Info("Getting events for month", zap.Int64("user_id", event.UserID))
	return s.repo.GetEventsForMonth(ctx, event)
}

func (s *CalendarService) CloseRepo() {
	s.log.Info("Closing repository")
	s.repo.Close()
	s.log.Info("Repository closed")
}
