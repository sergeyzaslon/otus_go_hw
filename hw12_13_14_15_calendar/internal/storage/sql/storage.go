package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
)

type Storage struct {
	ctx  context.Context
	conn *pgx.Conn
	dsn  string
	logg app.Logger
}

func New(ctx context.Context, dsn string, logg app.Logger) *Storage {
	return &Storage{
		dsn:  dsn,
		ctx:  ctx,
		logg: logg,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, s.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	s.conn = conn

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) Create(e app.Event) error {
	sql := `INSERT INTO 
				events(id, title, date, duration, description, userId, notify_before, notify_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	s.logg.Debug("SQL.Create: %s", sql)

	_, err := s.conn.Exec(
		s.ctx,
		sql,
		e.ID.String(),
		e.Title,
		e.Dt.Format(time.RFC3339),
		e.Duration.Seconds(),
		e.Description,
		e.UserID,
		e.NotifyBefore.Seconds(),
		e.Dt.Add(-1*e.NotifyBefore).Format(time.RFC3339),
	)

	return err
}

func (s *Storage) Update(e app.Event) error {
	sql := `UPDATE events 
				SET title=$1, date = $2, duration=$3, description=$4, userId=$5, notify_before=$6, notify_at=$7 
				WHERE id = $8`
	_, err := s.conn.Exec(
		s.ctx,
		sql,
		e.Title,
		e.Dt.Format(time.RFC3339),
		e.Duration.Seconds(),
		e.Description,
		e.UserID,
		e.NotifyBefore.Seconds(),
		e.Dt.Add(-1*e.NotifyBefore).Format(time.RFC3339),
		e.ID.String(),
	)

	return err
}

func (s *Storage) Delete(id uuid.UUID) error {
	sql := "DELETE FROM events WHERE id = $1"
	_, err := s.conn.Exec(s.ctx, sql, id)

	return err
}

func (s *Storage) FindOne(id uuid.UUID) (*app.Event, error) {
	var e app.Event
	var durationSeconds, notifyBeforeSeconds int

	query := "SELECT id, title, date, duration, description, userId, notify_before FROM events WHERE id = $1"
	err := s.conn.QueryRow(s.ctx, query, id).Scan(
		&e.ID,
		&e.Title,
		&e.Dt,
		&durationSeconds,
		&e.Description,
		&e.UserID,
		&notifyBeforeSeconds,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to scan SQL result into struct: %w", err)
	}

	e.Duration = time.Duration(durationSeconds * 1000000000)
	e.NotifyBefore = time.Duration(notifyBeforeSeconds * 1000000000)

	return &e, nil
}

func (s *Storage) FindAll() ([]app.Event, error) {
	sql := "SELECT id, title, date, duration, description, userId, notify_before FROM events ORDER BY date"
	rows, err := s.conn.Query(s.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return createEventsFromRows(rows)
}

func (s *Storage) GetEventsReadyToNotify(dt time.Time) ([]app.Event, error) {
	sql := `SELECT id, title, date, duration, description, userId, notify_before 
			FROM events 
			WHERE notified = false AND notify_at <= $1`
	rows, err := s.conn.Query(s.ctx, sql, dt.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return createEventsFromRows(rows)
}

func (s *Storage) GetEventsOlderThan(dt time.Time) ([]app.Event, error) {
	sql := `SELECT id, title, date, duration, description, userId, notify_before 
			FROM events 
			WHERE date <= $1`
	rows, err := s.conn.Query(s.ctx, sql, dt.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return createEventsFromRows(rows)
}

func createEventsFromRows(rows pgx.Rows) ([]app.Event, error) {
	var events []app.Event

	for rows.Next() {
		var e app.Event
		var durationSeconds, notifyBeforeSeconds int
		var description sql.NullString
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Dt,
			&durationSeconds,
			&description,
			&e.UserID,
			&notifyBeforeSeconds,
		); err != nil {
			return nil, fmt.Errorf("failed to scan SQL result into struct: %w", err)
		}

		if description.Valid {
			e.Description = description.String
		}

		e.Duration = time.Duration(durationSeconds * 1000000000)
		e.NotifyBefore = time.Duration(notifyBeforeSeconds * 1000000000)

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) MarkEventNotified(id uuid.UUID) error {
	sql := "UPDATE events SET notified = true WHERE id = $1"
	_, err := s.conn.Exec(
		s.ctx,
		sql,
		id.String(),
	)

	return err
}
