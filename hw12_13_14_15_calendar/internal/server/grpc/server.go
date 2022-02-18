//go:generate protoc --go_out=. --go-grpc_out=. ../../../api/EventService.proto --proto_path=../../../api

package internalgrpc

import (
	context "context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	application "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	grpc "google.golang.org/grpc"
)

type Server struct {
	UnimplementedEventServiceServer
	port    int
	grpcSrv *grpc.Server
	app     *application.App
	logg    application.Logger
}

func NewServerLogger(logger application.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		logger.Info("NEW GRPC Request: %v", req)
		return handler(ctx, req)
	}
}

func NewServer(port int, app *application.App, logger application.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			NewServerLogger(logger),
		),
	)

	s := &Server{
		port:    port,
		grpcSrv: grpcServer,
		app:     app,
		logg:    logger,
	}

	RegisterEventServiceServer(s.grpcSrv, s)

	return s
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", ":"+strconv.Itoa(s.port))
	if err != nil {
		return fmt.Errorf("failed to start GRPC service: %w", err)
	}

	s.logg.Info("Start GRPC Server on %d", s.port)

	return s.grpcSrv.Serve(lsn)
}

func (s *Server) Stop() {
	s.grpcSrv.GracefulStop()
}

func (s *Server) Create(ctx context.Context, in *Event) (*EventResponse, error) { // nolint:dupl
	appEvent := application.Event{
		Title:        in.GetTitle(),
		Duration:     time.Second * time.Duration(in.GetDuration()),
		Description:  in.GetDescription(),
		UserID:       in.GetUserId(),
		NotifyBefore: time.Second * time.Duration(in.GetNotifyBeforeSeconds()),
	}

	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid Id value. Exprected UUID, got %s, %w", in.GetId(), err)
	}
	appEvent.ID = id

	dt, err := time.Parse("2006-01-02 15:04:05", in.GetDate())
	if err != nil {
		return nil, fmt.Errorf("invalid date value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvent.Dt = dt

	if err = s.app.CreateEvent(ctx, appEvent); err != nil {
		return ResponseError(err.Error()), nil
	}

	return ResponseSuccess(), nil
}

func (s *Server) Update(ctx context.Context, in *Event) (*EventResponse, error) { // nolint:dupl
	appEvent := application.Event{
		Title:        in.GetTitle(),
		Duration:     time.Second * time.Duration(in.GetDuration()),
		Description:  in.GetDescription(),
		UserID:       in.GetUserId(),
		NotifyBefore: time.Second * time.Duration(in.GetNotifyBeforeSeconds()),
	}

	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid Id value. Exprected UUID, got %s, %w", in.GetId(), err)
	}
	appEvent.ID = id

	dt, err := time.Parse("2006-01-02 15:04:05", in.GetDate())
	if err != nil {
		return nil, fmt.Errorf("invalid date value. Exprected 2006-01-02 15:04:05, got %s,%w", in.GetId(), err)
	}
	appEvent.Dt = dt

	if err = s.app.UpdateEvent(ctx, appEvent); err != nil {
		return ResponseError(err.Error()), nil
	}

	return ResponseSuccess(), nil
}

func (s *Server) Delete(ctx context.Context, in *DeleteEventRequest) (*EventResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid Id value. Exprected UUID, got %s,%w", in.GetId(), err)
	}

	if err = s.app.DeleteEvent(ctx, id); err != nil {
		return ResponseError(err.Error()), nil
	}

	return ResponseSuccess(), nil
}

func (s *Server) EventListDay(ctx context.Context, in *EventListRequest) (*EventListResponse, error) {
	dt, err := time.Parse("2006-01-02", in.GetDate())
	if err != nil {
		return nil, fmt.Errorf("invalid date value. Expected yyyy-mm-dd, got %s", in.GetDate())
	}

	events, err := s.app.GetEventsByDay(ctx, dt)
	if err != nil {
		return nil, err
	}

	return ListResponseSuccess(events), nil
}

func (s *Server) EventListWeek(ctx context.Context, in *EventListRequest) (*EventListResponse, error) {
	dt, err := time.Parse("2006-01-02", in.GetDate())
	if err != nil {
		return nil, fmt.Errorf("invalid date value. Expected yyyy-mm-dd, got %s", in.GetDate())
	}

	events, err := s.app.GetEventsByWeek(ctx, dt)
	if err != nil {
		return nil, err
	}

	return ListResponseSuccess(events), nil
}

func (s *Server) EventListMonth(ctx context.Context, in *EventListRequest) (*EventListResponse, error) {
	dt, err := time.Parse("2006-01-02", in.GetDate())
	if err != nil {
		return nil, fmt.Errorf("invalid date value. Expected yyyy-mm-dd, got %s", in.GetDate())
	}

	events, err := s.app.GetEventsByMonth(ctx, dt)
	if err != nil {
		return nil, err
	}

	return ListResponseSuccess(events), nil
}

func ResponseSuccess() *EventResponse {
	return &EventResponse{
		Success: true,
		Error:   "",
	}
}

func ListResponseSuccess(events []application.Event) *EventListResponse {
	resp := EventListResponse{}
	for _, e := range events {
		resp.Events = append(resp.Events, &Event{
			Id:                  e.ID.String(),
			Title:               e.Title,
			Date:                e.Dt.Format(time.RFC3339),
			Duration:            uint32(e.Duration.Seconds()),
			Description:         e.Description,
			UserId:              e.UserID,
			NotifyBeforeSeconds: uint32(e.NotifyBefore.Seconds()),
		})
	}

	return &resp
}

func ResponseError(msg string) *EventResponse {
	return &EventResponse{
		Success: false,
		Error:   msg,
	}
}
