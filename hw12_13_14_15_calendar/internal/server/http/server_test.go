package internalhttp

import (
	"bytes"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestHttpServerHelloWorld(t *testing.T) {
	// Test Hello World
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Проверю тут сам роутинг + обработчики
	httpHandlers := NewRouter(createApp(t))
	httpHandlers.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, "Hello, world!\n", string(body))
}

func TestHttpServerEventsCrud(t *testing.T) {
	// Test Hello World
	body := bytes.NewBufferString(`{
		"id": "4927aa58-a175-429a-a125-c04765597152",
		"title": "Test Event 01",
		"description": "Test Event Description 01",
		"date": "2021-12-20 12:30:00",
		"duration": 60,
		"user_id": "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee",
		"notify_before_seconds": 60
	}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	httpHandlers := NewRouter(createApp(t))
	httpHandlers.ServeHTTP(w, req)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)
	respExpected := `{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Test Event 01","date":"2021-12-20 12:30:00","duration":60,"description":"Test Event Description 01","user_id":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notify_before_seconds":60}` // nolint:lll
	require.Equal(t, respExpected, string(respBody))

	// Прочитаем то, что создали
	req = httptest.NewRequest("GET", "/events", nil)
	w = httptest.NewRecorder()
	httpHandlers.ServeHTTP(w, req)
	resp = w.Result()
	respBody, _ = io.ReadAll(resp.Body)
	respExpected = `[{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Test Event 01","date":"2021-12-20 12:30:00","duration":60,"description":"Test Event Description 01","user_id":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notify_before_seconds":60}]` // nolint:lll
	require.Equal(t, respExpected, string(respBody))

	// Обновим:
	// Test Hello World
	body = bytes.NewBufferString(`{
		"title": "Test Event 01 UPD",
		"description": "Test Event Description 01 UPD",
		"date": "2021-12-21 12:30:00",
		"duration": 120,
		"user_id": "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee",
		"notify_before_seconds": 120
	}`)
	req = httptest.NewRequest("PUT", "/events/4927aa58-a175-429a-a125-c04765597152", body)
	w = httptest.NewRecorder()

	httpHandlers.ServeHTTP(w, req)

	resp = w.Result()
	respBody, _ = io.ReadAll(resp.Body)
	respExpected = `{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Test Event 01 UPD","date":"2021-12-21 12:30:00","duration":120,"description":"Test Event Description 01 UPD","user_id":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notify_before_seconds":120}` // nolint:lll
	require.Equal(t, respExpected, string(respBody))

	// Прочитаем то, что создали
	req = httptest.NewRequest("GET", "/events", nil)
	w = httptest.NewRecorder()
	httpHandlers.ServeHTTP(w, req)
	resp = w.Result()
	respBody, _ = io.ReadAll(resp.Body)
	respExpected = `[{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Test Event 01 UPD","date":"2021-12-21 12:30:00","duration":120,"description":"Test Event Description 01 UPD","user_id":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notify_before_seconds":120}]` // nolint:lll
	require.Equal(t, respExpected, string(respBody))
}

func createApp(t *testing.T) (*app.App, *logger.Logger) {
	t.Helper()
	logFile, err := os.CreateTemp("", "log")
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	logger, err := logger.New(logFile.Name(), "debug", "text_simple")
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	inMemoryStorage := memory.New()

	return app.New(logger, inMemoryStorage), logger
}
