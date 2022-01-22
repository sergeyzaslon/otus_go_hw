package sql

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/tools"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const AppConfFile = "configs/config.yaml"

func TestStorage(t *testing.T) {
	// Этот тест чисто для себя, в ci гонять его не надо:
	if _, err := os.Stat(AppConfFile); errors.Is(err, os.ErrNotExist) {
		t.Skip(AppConfFile + " file does not exists")
	}

	cfgContent, _ := os.ReadFile(AppConfFile)
	var cfg struct {
		Storage struct {
			Dsn string
		}
	}

	yaml.Unmarshal(cfgContent, cfg)

	ctx := context.Background()

	logg, _ := logger.New("stderr", "info", "text")
	s := New(ctx, cfg.Storage.Dsn, logg)
	if err := s.Connect(ctx); err != nil {
		t.Fatal("Failed to connect to DB server", err)
	}

	t.Run("test SQLStorage CRUDL", func(t *testing.T) {
		tx, _ := s.conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:       pgx.Serializable,
			AccessMode:     pgx.ReadWrite,
			DeferrableMode: pgx.NotDeferrable,
		})

		userID := "3b6394b3-acc6-4fd5-b8ce-3cbdf30745ef"
		dt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 12:00:00")

		event := app.NewEvent("Test Event", dt, time.Minute*30, userID)
		event.Description = "OTUS GoLang Lesson"
		event.NotifyBefore = time.Minute * 15

		err := s.Create(*event)
		require.Nil(t, err)

		saved, err := s.FindAll()
		require.Nil(t, err)
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
		err = s.Update(*event)
		if err != nil {
			t.Fatalf("Update failed: %s", err)
		}

		// Теперь он должен быть изменён
		saved, err = s.FindAll()
		if err != nil {
			t.Fatalf("failed to findAll after update: %s", err)
		}
		require.Len(t, saved, 1)
		require.Equal(t, *event, saved[0])

		// Удаляем объект
		s.Delete(event.ID)

		saved, _ = s.FindAll()
		require.Len(t, saved, 0)

		tx.Rollback(ctx)
	})

	t.Run("Test Storage::GetEventsReadyToNotify()", func(t *testing.T) {
		tx, _ := s.conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:       pgx.Serializable,
			AccessMode:     pgx.ReadWrite,
			DeferrableMode: pgx.NotDeferrable,
		})

		events := []app.Event{
			{
				ID:           tools.CreateUUID("4927aa58-a175-429a-a125-c04765597150"),
				Dt:           tools.CreateDate("2022-01-06T11:59:59Z"),
				NotifyBefore: 0,
			},
			{
				ID:           tools.CreateUUID("4927aa58-a175-429a-a125-c04765597151"),
				Dt:           tools.CreateDate("2022-01-06T12:00:00Z"),
				NotifyBefore: 0,
			},
			{
				ID:           tools.CreateUUID("4927aa58-a175-429a-a125-c04765597152"),
				Dt:           tools.CreateDate("2022-01-06T12:15:00Z"),
				NotifyBefore: time.Minute * 15,
			},
			{
				ID:           tools.CreateUUID("4927aa58-a175-429a-a125-c04765597153"),
				Dt:           tools.CreateDate("2022-01-06T12:00:01Z"),
				NotifyBefore: 0,
			},
		}

		for _, e := range events {
			s.Create(e)
		}

		readyEvents, err := s.GetEventsReadyToNotify(tools.CreateDate("2022-01-06T12:00:00Z"))
		require.Nil(t, err)

		ids := tools.ExtractEventID(readyEvents)
		idsExpected := []string{
			"4927aa58-a175-429a-a125-c04765597150",
			"4927aa58-a175-429a-a125-c04765597151",
			"4927aa58-a175-429a-a125-c04765597152",
		}
		require.Equal(t, idsExpected, ids)

		tx.Rollback(ctx)
	})
}
