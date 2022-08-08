package memorystorage

import (
	"errors"
	"testing"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

type LoggerMok struct{}

func (l LoggerMok) Error(msg string) {}

func TestStorage(t *testing.T) {
	s := New(LoggerMok{})

	t.Run("create event", func(t *testing.T) {
		e := storage.Event{ID: "test", Title: "test title"}
		s.CreateEvent(e)

		eInStorage, _ := s.GetEventByID(e.ID)

		require.Equal(t, e, eInStorage)
	})

	t.Run("get not found", func(t *testing.T) {
		_, err := s.GetEventByID("new test")

		require.True(t, errors.Is(storage.ErrEventNotFound, err))
	})

	t.Run("delete", func(t *testing.T) {
		e := storage.Event{ID: "test1", Title: "test title"}
		s.CreateEvent(e)
		err := s.DeleteEvent(e)

		require.Equal(t, nil, err)
	})

	t.Run("get all", func(t *testing.T) {
		e := storage.Event{ID: "test", Title: "test title"}

		events, _ := s.GetAllEvents()

		require.Equal(t, events, []storage.Event{e})
	})
}
