package internalgrpc

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/api/pb"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	logg := logger.New("ERROR", os.Stderr, os.Stdout)
	memoryStorage := memorystorage.New(logg)
	calendar := app.New(logg, memoryStorage)

	pb.RegisterEventsServer(s, NewService(calendar))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestApi(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	client := pb.NewEventsClient(conn)

	appEvent := storage.Event{
		Title:       "test",
		Description: "test",
		StartAt:     time.Now(),
		EndAt:       time.Now(),
		RemindFor:   time.Now(),
		UserID:      123,
		ID:          "",
	}

	t.Run("Create Event", func(t *testing.T) {
		resp, err := client.CreateEvent(ctx, &pb.CreateEventRequest{
			Event: convertAppEventToPbEvent(appEvent),
		})
		if err != nil {
			t.Fatalf("CreateEvent failed: %v", err)
		}

		appEvent.ID = resp.Event.Id

		require.Equal(t, nil, err)
	})

	t.Run("Update Event", func(t *testing.T) {
		appEvent.Title = "new title"
		resp, err := client.UpdateEvent(ctx, &pb.UpdateEventRequest{
			Event: convertAppEventToPbEvent(appEvent),
		})
		if err != nil {
			t.Fatalf("Update event failed: %v", err)
		}

		require.Equal(t, appEvent.Title, convertPbEventToAppEvent(resp.Event).Title)
	})
}
