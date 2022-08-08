package memorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
	log    Logger
}
type Logger interface {
	Error(msg string)
}

func New(l Logger) *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
		log:    l,
	}
}

func (s *Storage) CreateEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[e.ID] = e

	return nil
}

func (s *Storage) DeleteEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, e.ID)

	return nil
}

func (s *Storage) Update(e storage.Event) error {
	s.events[e.ID] = e

	return nil
}

func (s *Storage) GetEventByID(id string) (storage.Event, error) {
	e, ok := s.events[id]

	if !ok {
		s.log.Error(fmt.Sprintf(storage.ErrEventNotFound.Error(), e.ID))
		return storage.Event{}, storage.ErrEventNotFound
	}

	return e, nil
}

func (s *Storage) GetAllEvents() ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	for _, event := range s.events {
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) MigrationUp(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}
