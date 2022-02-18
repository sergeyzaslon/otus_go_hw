package app

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDateBusy      = errors.New("date is already occupied")
	ErrRequiredField = errors.New("required field")
)

type Event struct {
	ID           uuid.UUID
	Title        string
	Dt           time.Time
	Duration     time.Duration
	Description  string
	UserID       string
	NotifyBefore time.Duration
}

func NewEvent(title string, dt time.Time, duration time.Duration, userID string) *Event {
	id, _ := uuid.NewRandom()
	return &Event{
		ID:       id,
		Title:    title,
		Dt:       dt,
		Duration: duration,
		UserID:   userID,
	}
}

type Notification struct {
	EventID uuid.UUID
	Title   string
	Dt      time.Time
	UserID  string
}
