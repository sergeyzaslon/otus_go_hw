package memory

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
)

type Storage struct {
	mu       sync.RWMutex
	events   map[uuid.UUID]app.Event
	notified map[uuid.UUID]bool
}

func New() *Storage {
	return &Storage{
		events:   make(map[uuid.UUID]app.Event),
		notified: make(map[uuid.UUID]bool),
	}
}

func (s *Storage) Create(e app.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[e.ID] = e

	return nil
}

func (s *Storage) Update(e app.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[e.ID] = e

	return nil
}

func (s *Storage) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, id)

	return nil
}

func (s *Storage) FindOne(id uuid.UUID) (*app.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if event, ok := s.events[id]; ok {
		return &event, nil
	}

	return nil, nil
}

func (s *Storage) FindAll() ([]app.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]app.Event, 0, len(s.events))
	for _, e := range s.events {
		events = append(events, e)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Dt.Unix() < events[j].Dt.Unix()
	})

	return events, nil
}

func (s *Storage) GetEventsReadyToNotify(dt time.Time) ([]app.Event, error) {
	var res []app.Event

	for _, e := range s.events {
		if e.Dt.Sub(dt) <= e.NotifyBefore && !s.notified[e.ID] {
			res = append(res, e)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Dt.Unix() < res[j].Dt.Unix()
	})

	return res, nil
}

func (s *Storage) GetEventsOlderThan(dt time.Time) ([]app.Event, error) {
	var res []app.Event

	for _, e := range s.events {
		if dt.Sub(e.Dt) >= 0 {
			res = append(res, e)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Dt.Unix() < res[j].Dt.Unix()
	})

	return res, nil
}

func (s *Storage) MarkEventNotified(id uuid.UUID) error {
	s.notified[id] = true
	return nil
}
