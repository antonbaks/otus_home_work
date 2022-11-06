package integration_test

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/api/pb"
	internalgrpc "github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const UserID = 123

type IntegrationSuite struct {
	suite.Suite
	c       pb.EventsClient
	ctx     context.Context
	eventID string
}

func (s *IntegrationSuite) SetupSuite() {
	flag.Parse()

	conn, err := grpc.Dial("calendar:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(s.T(), err)
	s.c = pb.NewEventsClient(conn)
	s.ctx = context.Background()
}

func (s *IntegrationSuite) TestCreateEvent() {
	timeNow := time.Now()
	appEvent := storage.Event{
		Title:       "test",
		Description: "test",
		StartAt:     timeNow.Add(time.Hour),
		EndAt:       timeNow.Add(time.Hour * 48),
		UserID:      UserID,
		RemindFor:   timeNow,
	}

	resp, err := s.c.CreateEvent(s.ctx, &pb.CreateEventRequest{
		Event: internalgrpc.ConvertAppEventToPbEvent(appEvent),
	})
	s.Require().NoError(err)

	s.eventID = resp.Event.Id
}

func (s *IntegrationSuite) TestListEvent() {
	resp, err := s.c.GetEvents(s.ctx, &pb.GetEventsRequest{
		UserId:  UserID,
		StartAt: timestamppb.New(time.Now().Add(time.Hour * -1)),
		EndAt:   timestamppb.New(time.Now().Add(time.Hour)),
	})

	s.Require().NoError(err)

	s.Require().Equal(resp.Events[0].Id, s.eventID)
}

func (s *IntegrationSuite) TestSendNotification() {
	time.Sleep(time.Minute)

	resp, err := s.c.IsNotifyEvent(s.ctx, &pb.IsNotifyEventRequest{
		UserId:  UserID,
		EventId: s.eventID,
	})

	s.Require().NoError(err)

	s.Require().Equal(resp.Event.Id, s.eventID)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
