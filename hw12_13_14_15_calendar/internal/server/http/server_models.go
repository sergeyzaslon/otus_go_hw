package internalhttp

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
)

type EventDto struct {
	ID                  string `json:"id"`
	Title               string `json:"title"`
	Date                string `json:"date"`
	Duration            uint32 `json:"duration"`
	Description         string `json:"description"`
	UserID              string `json:"userId"`
	NotifyBeforeSeconds uint32 `json:"notifyBeforeSeconds"`
}

type ErrorDto struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (e *EventDto) GetModel() (*app.Event, error) {
	dt, err := time.Parse("2006-01-02 15:04:00", e.Date)
	if err != nil {
		return nil, fmt.Errorf("date exprected to be 'yyyy-mm-dd hh:mm:ss', got: %s, %w", e.Date, err)
	}

	duration := time.Second & time.Duration(e.Duration)

	notifyBefore := time.Second * time.Duration(e.NotifyBeforeSeconds)

	id, err := uuid.Parse(e.ID)
	if err != nil {
		return nil, fmt.Errorf("id exprected to be uuid, got: %s, %w", e.ID, err)
	}

	appEvent := app.NewEvent(e.Title, dt, duration, e.UserID)
	appEvent.Description = e.Description
	appEvent.ID = id
	appEvent.NotifyBefore = notifyBefore

	return appEvent, nil
}

func CreateEventDtoFromModel(event app.Event) EventDto {
	dto := EventDto{}
	dto.ID = event.ID.String()
	dto.Title = event.Title
	dto.Date = event.Dt.Format(time.RFC3339)
	dto.Duration = uint32(event.Duration.Seconds())
	dto.Description = event.Description
	dto.UserID = event.UserID
	dto.NotifyBeforeSeconds = uint32(event.NotifyBefore.Seconds())

	return dto
}
