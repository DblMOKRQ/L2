package repository

import (
	"awesomeProject/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func (s *Storage) NewRepository() *Repository {
	return &Repository{db: s.db, log: s.log.Named("repository")}
}

const (
	createQuery = `
		INSERT INTO calendar (user_id, date,event) VALUES ($1,$2,$3) `
	updateQuery     = `UPDATE calendar SET date = $1, event = $2 WHERE user_id = $3`
	deleteQuery     = `DELETE FROM calendar WHERE user_id = $1 AND date = $2 `
	getForDayQuery  = `SELECT user_id,date,event FROM calendar WHERE user_id = $1 AND date = $2`
	getForWeekQuery = `SELECT user_id,date,event FROM calendar WHERE user_id = $1 
                                    AND date >= $2::date 
                                    AND date < $2::date + INTERVAL '7 day' 
                                ORDER BY date;`
	getForMouthQuery = `    SELECT user_id,date, event 
    FROM calendar 
    WHERE user_id = $1 
      AND date >= $2::date 
      AND date <= $2::date + INTERVAL '1 month'
    ORDER BY date;`
)

func (r *Repository) CreateEvent(ctx context.Context, event *models.Event) error {
	r.log.Debug("Creating Event", zap.Any("event", event))
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.log.Error("Error begin transaction", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, createQuery,
		event.UserID,
		event.Date,
		event.Event,
	)
	if err != nil {
		r.log.Error("Error create event", zap.Error(err))
		return fmt.Errorf("failed to create event: %w", err)
	}
	r.log.Debug("Created event", zap.Any("event", event))

	return tx.Commit(ctx)
}
func (r *Repository) UpdateEvent(ctx context.Context, event *models.Event) error {
	r.log.Debug("Updating Event", zap.Any("event", event))
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.log.Error("Error begin transaction", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	_, err = tx.Exec(ctx, updateQuery,
		event.Date,
		event.Event,
		event.UserID,
	)
	if err != nil {
		r.log.Error("Error update event", zap.Error(err))
		return fmt.Errorf("failed to update event: %w", err)
	}
	r.log.Debug("Updated event", zap.Any("event", event))
	return tx.Commit(ctx)
}
func (r *Repository) DeleteEvent(ctx context.Context, event *models.Event) error {
	r.log.Debug("Deleting Event", zap.Any("event", event))
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.log.Error("Error begin transaction", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	_, err = tx.Exec(ctx, deleteQuery,
		event.UserID,
		event.Date,
	)
	if err != nil {
		r.log.Error("Error delete event", zap.Error(err))
		return fmt.Errorf("failed to delete event: %w", err)
	}
	r.log.Debug("Deleted event", zap.Any("event", event))
	return tx.Commit(ctx)
}
func (r *Repository) GetEventsForDay(ctx context.Context, event *models.Event) ([]models.Event, error) {
	r.log.Debug("Getting Events for Day", zap.Any("event", event))
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.log.Error("Error begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	queryEvents, err := tx.Query(ctx, getForDayQuery, event.UserID, event.Date)
	if err != nil {
		r.log.Error("Error get events for day", zap.Error(err))
		return nil, fmt.Errorf("failed to get events for day: %w", err)
	}
	var events []models.Event
	for queryEvents.Next() {
		var ev models.Event
		err = queryEvents.Scan(&ev.UserID, &ev.Date, &ev.Event)
		if err != nil {
			r.log.Error("Error get events for day", zap.Error(err))
			return nil, fmt.Errorf("failed to get events for day: %w", err)
		}
		events = append(events, ev)
	}
	queryEvents.Close()
	if err := queryEvents.Err(); err != nil {
		r.log.Error("Error get events for day", zap.Error(err))
		return nil, fmt.Errorf("failed to get events for day: %w", err)
	}
	r.log.Debug("Got events for day", zap.Int("events", len(events)))
	return events, nil
}
func (r *Repository) GetEventsForWeek(ctx context.Context, event *models.Event) ([]models.Event, error) {
	r.log.Debug("Getting Events for Week", zap.Any("event", event))
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.log.Error("Error begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	queryEvents, err := tx.Query(ctx, getForWeekQuery, event.UserID, event.Date)
	if err != nil {
		r.log.Error("Error get events for week", zap.Error(err))
		return nil, fmt.Errorf("failed to get events for week: %w", err)
	}
	var events []models.Event
	for queryEvents.Next() {
		var ev models.Event
		err = queryEvents.Scan(&ev.UserID, &ev.Date, &ev.Event)
		if err != nil {
			r.log.Error("Error get events for week", zap.Error(err))
			return nil, fmt.Errorf("failed to get events for week: %w", err)
		}
		events = append(events, ev)
	}
	queryEvents.Close()
	if err := queryEvents.Err(); err != nil {
		r.log.Error("Error get events for week", zap.Error(err))
		return nil, fmt.Errorf("failed to get events for week: %w", err)
	}
	return events, nil
}
func (r *Repository) GetEventsForMonth(ctx context.Context, event *models.Event) ([]models.Event, error) {
	r.log.Debug("Getting Events for Month", zap.Any("event", event))
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.log.Error("Error begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	queryEvents, err := tx.Query(ctx, getForMouthQuery, event.UserID, event.Date)
	if err != nil {
		r.log.Error("Error get events for month", zap.Error(err))
		return nil, fmt.Errorf("failed to get events for month: %w", err)
	}
	var events []models.Event
	for queryEvents.Next() {
		var ev models.Event
		err = queryEvents.Scan(&ev.UserID, &ev.Date, &ev.Event)
		if err != nil {
			r.log.Error("Error get events for mouth", zap.Error(err))
			return nil, fmt.Errorf("failed to get events for month: %w", err)
		}
		events = append(events, ev)
	}
	queryEvents.Close()
	if err := queryEvents.Err(); err != nil {
		r.log.Error("Error get events for month", zap.Error(err))
		return nil, fmt.Errorf("failed to get events for month: %w", err)
	}
	return events, nil
}

func (r *Repository) Close() {
	r.log.Info("Closing repository")
	r.db.Close()
}
