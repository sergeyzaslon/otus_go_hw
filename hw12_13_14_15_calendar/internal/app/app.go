package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type App struct {
	logg Logger
	repo Storage
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logg: logger,
		repo: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, evt Event) error {
	var (
		err  error
		prev *Event
	)
	a.logg.Debug("App.CreateEvent - create new event")

	// Проверим, что он уже существует
	if prev, err = a.repo.FindOne(evt.ID); err != nil {
		a.logg.Error("App.CreateEvent ERROR: %s", err)
		return err
	}

	if prev != nil {
		a.logg.Warn("App.CreateEvent.AlreadyExists: %s", evt.ID)
		return fmt.Errorf("validation error: event with such id already exists: %s", evt.ID)
	}

	// Если ещё нет с таким ID - создаём
	if err = a.repo.Create(evt); err != nil {
		a.logg.Error("App.CreateEvent ERROR: %s", err)
		return err
	}

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, evt Event) error {
	var (
		err  error
		prev *Event
	)
	a.logg.Debug("App.UpdateEvent.Begin %s", evt.ID)

	// Проверим наличие
	if prev, err = a.repo.FindOne(evt.ID); err != nil {
		a.logg.Error("App.UpdateEvent ERROR: %s", err)
		return err
	}

	if prev == nil {
		a.logg.Warn("App.UpdateEvent.NotFound: %s", evt.ID)
		return fmt.Errorf("validation error: event %s not found", evt.ID)
	}

	if err = a.repo.Update(evt); err != nil {
		a.logg.Error("App.UpdateEvent ERROR: %s", err)
		return err
	}

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	var (
		err  error
		prev *Event
	)
	a.logg.Debug("App.DeleteEvent.Begin %s", id)

	// Проверим наличие
	if prev, err = a.repo.FindOne(id); err != nil {
		a.logg.Error("App.DeleteEvent ERROR: %s", err)
		return err
	}

	if prev == nil {
		a.logg.Warn("App.DeleteEvent.NotFound: %s", id)
		return fmt.Errorf("validation error: event %s not found", id)
	}

	// Если ещё нет с таким ID - создаём
	if err = a.repo.Delete(prev.ID); err != nil {
		a.logg.Error("App.DeleteEvent ERROR: %s", err)
		return err
	}

	return nil
}

func (a *App) GetEvents(ctx context.Context) ([]Event, error) {
	return a.repo.FindAll()
}

func (a *App) GetEventsByDay(ctx context.Context, day time.Time) ([]Event, error) {
	finish := day.AddDate(0, 0, 1)
	return a.GetEventsByInterval(ctx, day, finish.Sub(day))
}

func (a *App) GetEventsByWeek(ctx context.Context, day time.Time) ([]Event, error) {
	finish := day.AddDate(0, 0, 7)
	return a.GetEventsByInterval(ctx, day, finish.Sub(day))
}

func (a *App) GetEventsByMonth(ctx context.Context, day time.Time) ([]Event, error) {
	finish := day.AddDate(0, 1, 0)
	return a.GetEventsByInterval(ctx, day, finish.Sub(day))
}

func (a *App) GetEventsByInterval(ctx context.Context, day time.Time, interval time.Duration) ([]Event, error) {
	events := make([]Event, 0)

	// 2022-01-02 12:45:13 ==> 2022-01-02 00:00:00
	day = day.Truncate(time.Minute * 1440)

	a.logg.Debug("Get Event List from %s, interval: %s", day, interval)

	items, err := a.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		diff := item.Dt.Sub(day)
		if diff >= 0 && diff < interval {
			fmt.Printf("%s + %s >= %s\n", day, interval, item.Dt)
			events = append(events, item)
		} else {
			fmt.Printf("%s + %s < %s\n", day, interval, item.Dt)
		}
	}

	return events, nil
}
