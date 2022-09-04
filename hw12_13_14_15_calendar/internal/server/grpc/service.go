package internalgrpc

import (
	"context"
	"errors"
	"time"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/api/pb"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrEmptyData = errors.New("request data is empty")

type App interface {
	CreateEvent(event storage.Event) (storage.Event, error)
	UpdateEvent(event storage.Event) (storage.Event, error)
	DeleteEvent(event storage.Event) error
	GetEvents(startAt time.Time, endAt time.Time, userID int) ([]storage.Event, error)
	IsNotify(id string, userID int) (storage.Event, error)
}

type Service struct {
	pb.UnsafeEventsServer
	app App
}

func NewService(app App) *Service {
	return &Service{
		app: app,
	}
}

func (s *Service) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	if req.Event == nil {
		return nil, ErrEmptyData
	}

	e, err := s.app.CreateEvent(convertPbEventToAppEvent(req.Event))
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{Event: ConvertAppEventToPbEvent(e)}, nil
}

func (s *Service) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	if err := s.app.DeleteEvent(storage.Event{ID: req.Id}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	if req.Event == nil {
		return nil, ErrEmptyData
	}

	e, err := s.app.UpdateEvent(convertPbEventToAppEvent(req.Event))
	if err != nil {
		return nil, err
	}

	return &pb.UpdateEventResponse{Event: ConvertAppEventToPbEvent(e)}, nil
}

func (s *Service) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	e, err := s.app.GetEvents(req.StartAt.AsTime(), req.EndAt.AsTime(), int(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.GetEventsResponse{Events: ConvertSliceAppEventToPbEvent(e)}, nil
}

func (s *Service) IsNotifyEvent(ctx context.Context, req *pb.IsNotifyEventRequest) (*pb.IsNotifyEventResponse, error) {
	e, err := s.app.IsNotify(req.EventId, int(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.IsNotifyEventResponse{Event: ConvertAppEventToPbEvent(e)}, nil
}

func convertPbEventToAppEvent(pbe *pb.Event) storage.Event {
	return storage.Event{
		ID:          pbe.Id,
		Title:       pbe.Title,
		Description: pbe.Description,
		StartAt:     pbe.StartAt.AsTime(),
		EndAt:       pbe.EndAt.AsTime(),
		UserID:      int(pbe.UserId),
		RemindFor:   pbe.RemindFor.AsTime(),
	}
}

func ConvertAppEventToPbEvent(e storage.Event) *pb.Event {
	return &pb.Event{
		Id:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		StartAt:     timestamppb.New(e.StartAt),
		EndAt:       timestamppb.New(e.EndAt),
		UserId:      uint32(e.UserID),
		RemindFor:   timestamppb.New(e.RemindFor),
	}
}

func ConvertSliceAppEventToPbEvent(appE []storage.Event) []*pb.Event {
	pbEvents := make([]*pb.Event, 0)
	for _, e := range appE {
		pbEvents = append(pbEvents, ConvertAppEventToPbEvent(e))
	}

	return pbEvents
}
