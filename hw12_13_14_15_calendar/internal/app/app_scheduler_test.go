package app_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/tools"
	"github.com/stretchr/testify/require"
)

type TstEventSource struct {
	storage *memory.Storage
}

func (s *TstEventSource) Add(id uuid.UUID, dt time.Time, nfy time.Duration) {
	e := app.NewEvent("Test", dt, time.Minute*30, "")
	e.ID = id
	e.NotifyBefore = nfy

	s.storage.Create(*e)
}

func (s *TstEventSource) MarkEventNotified(id uuid.UUID) error {
	return nil
}

func (s *TstEventSource) GetEventsReadyToNotify(dt time.Time) ([]app.Event, error) {
	return s.storage.GetEventsReadyToNotify(dt)
}

func (s *TstEventSource) GetEventsOlderThan(dt time.Time) ([]app.Event, error) {
	return s.storage.GetEventsOlderThan(dt)
}

func (s *TstEventSource) Delete(id uuid.UUID) error {
	return s.storage.Delete(id)
}

type TstNotificationReceiver struct {
	Notifications []app.Notification
}

func (r *TstNotificationReceiver) Add(n app.Notification) error {
	r.Notifications = append(r.Notifications, n)
	return nil
}

type TstClock struct {
	now time.Time
}

func (c *TstClock) NowIs(t time.Time) {
	c.now = t
}

func (c *TstClock) Now() time.Time {
	return c.now
}

func TestAppScheduler(t *testing.T) {
	es := &TstEventSource{memory.New()}
	rcv := &TstNotificationReceiver{}
	clck := &TstClock{}
	log, _ := logger.New("stderr", "error", "text")

	sheduler := app.NewAppScheduler(es, rcv, clck, log)

	clck.NowIs(tools.CreateDate("2022-01-06T12:00:00Z"))

	es.Add(tools.CreateUUID("4927aa58-a175-429a-a125-c04765597150"), tools.CreateDate("2022-01-06T11:00:00Z"), 0)              // nolint:lll
	es.Add(tools.CreateUUID("4927aa58-a175-429a-a125-c04765597151"), tools.CreateDate("2022-01-06T12:00:00Z"), 0)              // nolint:lll
	es.Add(tools.CreateUUID("4927aa58-a175-429a-a125-c04765597152"), tools.CreateDate("2022-01-06T12:00:01Z"), 0)              // nolint:lll
	es.Add(tools.CreateUUID("4927aa58-a175-429a-a125-c04765597153"), tools.CreateDate("2022-01-06T12:15:00Z"), time.Minute*15) // nolint:lll
	es.Add(tools.CreateUUID("4927aa58-a175-429a-a125-c04765597154"), tools.CreateDate("2022-01-06T12:16:00Z"), time.Minute*15) // nolint:lll
	es.Add(tools.CreateUUID("4927aa58-a175-429a-a125-c04765597155"), tools.CreateDate("2022-01-06T13:00:00Z"), time.Minute*15) // nolint:lll

	sheduler.Notify()

	var ids []string // nolint:prealloc
	for _, n := range rcv.Notifications {
		ids = append(ids, n.EventID.String())
	}

	expectedEventIDs := []string{
		"4927aa58-a175-429a-a125-c04765597150",
		"4927aa58-a175-429a-a125-c04765597151",
		"4927aa58-a175-429a-a125-c04765597153",
	}
	require.Equal(t, expectedEventIDs, ids)
}

func TestHowAppSchedulerRemovesOldEvents(t *testing.T) {
	storage := memory.New()
	es := &TstEventSource{storage}
	rcv := &TstNotificationReceiver{}
	clck := &TstClock{}
	log, _ := logger.New("stderr", "error", "text")

	sheduler := app.NewAppScheduler(es, rcv, clck, log)

	clck.NowIs(tools.CreateDate("2022-01-06T12:00:00Z"))

	es.Add(tools.CreateUUID("4927aa58-a175-429a-a125-c04765597150"), tools.CreateDate("2020-01-06T12:00:00Z"), 0) // nolint:lll
	es.Add(tools.CreateUUID("4927aa58-a175-429a-a125-c04765597151"), tools.CreateDate("2021-01-06T12:00:00Z"), 0) // nolint:lll
	es.Add(tools.CreateUUID("4927aa58-a175-429a-a125-c04765597152"), tools.CreateDate("2021-01-06T12:00:01Z"), 0) // nolint:lll

	sheduler.RemoveOneYearOld()

	events, err := storage.FindAll()
	if err != nil {
		t.Errorf("failed to get events: %s", err)
	}

	require.Len(t, events, 1)
	require.Equal(t, "4927aa58-a175-429a-a125-c04765597152", events[0].ID.String())
}
