package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	host   string
	port   int
	logger Logger
	server *http.Server
}

type Logger interface {
	Info(msg string, params ...interface{})
	LogHTTPRequest(r *http.Request, code, length int)
}

func NewServer(logger Logger, app *app.App, host string, port int) *Server {
	myServer := &Server{
		host:   host,
		port:   port,
		logger: logger,
		server: nil,
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, strconv.Itoa(port)),
		Handler: loggingMiddleware(NewRouter(app), logger),
	}

	myServer.server = httpServer

	return myServer
}

func NewRouter(app *app.App) http.Handler {
	handlers := NewServerHandlers(app)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HelloWorld).Methods("GET")

	return r
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Start HTTP Server on %s:%d", s.host, s.port)
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.Shutdown(ctx)
	return nil
}
