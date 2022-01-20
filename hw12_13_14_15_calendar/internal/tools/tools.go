package tools

import (
	"time"

	"github.com/google/uuid"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
)

func CreateUUID(str string) uuid.UUID {
	id, _ := uuid.Parse(str)
	return id
}

func CreateDate(str string) time.Time {
	dt, _ := time.Parse(time.RFC3339, str)
	return dt
}

func ExtractEventID(events []app.Event) []string {
	res := make([]string, 0, len(events))

	for _, e := range events {
		res = append(res, e.ID.String())
	}

	return res
}
