package memory

import (
	"testing"
	"time"

	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	s := New()
	t.Run("test inMemory storage CRUDL", func(t *testing.T) {
		userID := "3b6394b3-acc6-4fd5-b8ce-3cbdf30745ef"
		dt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 12:00:00")

		event := app.NewEvent("Test Event", dt, time.Minute*30, userID)
		event.Description = "OTUS GoLang Lesson"
		event.NotifyBefore = time.Minute * 15

		s.Create(*event)

		saved, _ := s.FindAll()
		require.Len(t, saved, 1)
		require.Equal(t, *event, saved[0])

		// Обновим параметры события:
		event.Title = "Test Event Upd"
		event.Description = "OTUS GoLang Lesson Upd"
		event.NotifyBefore = time.Minute * 15

		// Убедимся, что объект не изменился в репозитории только потому, что там хранятся ссылки,а не копии
		saved, _ = s.FindAll()
		require.Len(t, saved, 1)
		require.NotEqual(t, *event, saved[0])

		// Обновляем объект в репозитории
		s.Update(*event)

		// Теперь он должен быть изменён
		saved, _ = s.FindAll()
		require.Len(t, saved, 1)
		require.Equal(t, *event, saved[0])

		// Удаляем объект
		s.Delete(event.ID)

		saved, _ = s.FindAll()
		require.Len(t, saved, 0)
	})
}
