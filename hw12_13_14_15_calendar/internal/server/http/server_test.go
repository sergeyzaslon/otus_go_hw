package internalhttp

import (
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

func createApp(t *testing.T) *app.App {
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

	return app.New(logger, inMemoryStorage)
}
