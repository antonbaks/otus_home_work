package memorystorage

import (
	"fmt"
	"sync"
	"time"

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

func (s *Storage) GetEvents(startAt time.Time, endAt time.Time, userID int) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	for _, event := range s.events {
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) IsNotify(id string, userID int) (storage.Event, error) {
	e, err := s.GetEventByID(id)
	if err != nil {
		return e, err
	}

	if e.UserID == userID {
		return e, nil
	}

	return storage.Event{}, storage.ErrEventNotFound
}

func (s *Storage) MigrationUp() error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}
