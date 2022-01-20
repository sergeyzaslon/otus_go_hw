package internalhttp

import (
	"net/http"

	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
)

type ServerHandlers struct {
	app *app.App
}

func NewServerHandlers(a *app.App) *ServerHandlers {
	return &ServerHandlers{app: a}
}

func (s *ServerHandlers) HelloWorld(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Hello, world!\n")
	w.WriteHeader(200)
	w.Write(msg)
}
