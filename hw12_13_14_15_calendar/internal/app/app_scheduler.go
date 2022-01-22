package app

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventSource interface {
	GetEventsReadyToNotify(dt time.Time) ([]Event, error)
	GetEventsOlderThan(dt time.Time) ([]Event, error)
	Delete(id uuid.UUID) error
	MarkEventNotified(id uuid.UUID) error
}

type NotificationReceiver interface {
	Add(Notification) error
}

type Scheduler struct {
	eventSource          EventSource
	notificationReceiver NotificationReceiver
	clock                Clock
	logger               Logger
}

type Clock interface {
	Now() time.Time
}

func NewAppScheduler(es EventSource, rcv NotificationReceiver, clck Clock, logger Logger) *Scheduler {
	return &Scheduler{
		es,
		rcv,
		clck,
		logger,
	}
}

func (s *Scheduler) Notify() error {
	dt := s.clock.Now()
	events, err := s.eventSource.GetEventsReadyToNotify(dt)
	if err != nil {
		return fmt.Errorf("failed to get events for date %s: %w", dt, err)
	}

	if len(events) > 0 {
		s.logger.Info("[scheduler] Got %d messages to send", len(events))
	} else {
		s.logger.Debug("[scheduler] No messages to send")
	}

	for _, event := range events {
		notification := Notification{
			EventID: event.ID,
			Title:   event.Title,
			Dt:      event.Dt,
			UserID:  event.UserID,
		}
		if err := s.notificationReceiver.Add(notification); err != nil {
			return fmt.Errorf("failed to push notification for event %s:  %w", event.ID, err)
		}

		s.eventSource.MarkEventNotified(event.ID)

		s.logger.Info("[scheduler] Event %s notified", notification.EventID)
	}

	return nil
}

func (s *Scheduler) RemoveOneYearOld() error {
	dt := s.clock.Now()
	dtOneYearAgo := dt.AddDate(-1, 0, 0)

	events, err := s.eventSource.GetEventsOlderThan(dtOneYearAgo)
	if err != nil {
		return fmt.Errorf("failed to get events for date %s: %w", dt, err)
	}

	fmt.Println(dtOneYearAgo)
	for _, event := range events {
		fmt.Println(event.Dt)
		s.eventSource.Delete(event.ID)

		s.logger.Info("[scheduler] Old Event %s removed. Date was: %s", event.ID, event.Dt)
	}

	return nil
}
