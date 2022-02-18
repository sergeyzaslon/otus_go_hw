package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
)

type ServerHandlers struct {
	app *app.App
	log Logger
}

func NewServerHandlers(a *app.App, log Logger) *ServerHandlers {
	return &ServerHandlers{app: a, log: log}
}

func (s *ServerHandlers) HelloWorld(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Hello, world!\n")
	w.WriteHeader(200)
	w.Write(msg)
}

func (s *ServerHandlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var dto EventDto
	err := ParseRequest(r, &dto)
	if err != nil {
		s.RespondError(w, http.StatusBadRequest, err)
		return
	}

	event, err := dto.GetModel()
	if err != nil {
		s.RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.CreateEvent(r.Context(), *event)
	if err != nil {
		if errors.Is(err, app.ErrEventWithSuchIDAlreadyExists) {
			s.RespondError(w, http.StatusBadRequest, err)
		} else {
			s.RespondError(w, http.StatusInternalServerError, err)
		}
		return
	}

	responseData, err := json.Marshal(dto)
	if err != nil {
		s.RespondError(w, http.StatusInternalServerError, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

func (s *ServerHandlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var dto EventDto
	err := ParseRequest(r, &dto)
	if err != nil {
		s.RespondError(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	dto.ID = vars["id"]

	event, err := dto.GetModel()
	if err != nil {
		s.RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.UpdateEvent(r.Context(), *event)
	if err != nil {
		s.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	responseData, _ := json.Marshal(dto)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (s *ServerHandlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		s.RespondError(w, http.StatusBadRequest, err)
	}

	err = s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		s.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerHandlers) ListEvents(w http.ResponseWriter, r *http.Request) {
	dtStartStr := r.URL.Query().Get("date")
	dtIntervalStr := r.URL.Query().Get("interval")

	withDate := false
	dtStart, err := time.Parse("2006-01-02", dtStartStr)
	if err == nil {
		withDate = true
	}

	var events []app.Event
	if withDate {
		switch dtIntervalStr {
		case "day":
			events, err = s.app.GetEventsByDay(r.Context(), dtStart)
		case "week":
			events, err = s.app.GetEventsByWeek(r.Context(), dtStart)
		case "month":
			events, err = s.app.GetEventsByMonth(r.Context(), dtStart)
		default:
			events, err = s.app.GetEventsByDay(r.Context(), dtStart)
		}
	} else {
		events, err = s.app.GetEvents(r.Context())
	}

	if err != nil {
		s.RespondError(w, http.StatusInternalServerError, err)
	}

	eventDtos := make([]EventDto, 0, len(events))
	for _, e := range events {
		eventDtos = append(eventDtos, CreateEventDtoFromModel(e))
	}

	responseData, _ := json.Marshal(eventDtos)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func ParseRequest(r *http.Request, dto interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	err = json.Unmarshal(data, dto)
	if err != nil {
		return fmt.Errorf("failed to decode JSON request: %w", err)
	}

	return nil
}

func (s *ServerHandlers) RespondError(w http.ResponseWriter, code int, appError error) {
	data, err := json.Marshal(ErrorDto{
		false,
		appError.Error(),
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to marshall error dto"))
	}

	s.log.Error("HTTP ERR: %s", appError.Error())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
