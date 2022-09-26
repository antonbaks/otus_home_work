package app

import (
	"context"
	"time"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	log Logger
	s   Storage
}

type Logger interface {
	Error(msg string)
}

type Storage interface {
	CreateEvent(e storage.Event) error
	DeleteEvent(e storage.Event) error
	Update(e storage.Event) error
	GetEventByID(id string) (storage.Event, error)
	GetEvents(startAt time.Time, endAt time.Time, UserID int) ([]storage.Event, error)
	MigrationUp(ctx context.Context) error
	Close(ctx context.Context) error
}

func New(log Logger, s Storage) *App {
	return &App{
		log: log,
		s:   s,
	}
}

func (a *App) CreateEvent(event storage.Event) (storage.Event, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return storage.Event{}, err
	}

	event.ID = newUUID.String()

	if err := a.s.CreateEvent(event); err != nil {
		return storage.Event{}, err
	}

	return event, nil
}

func (a *App) UpdateEvent(event storage.Event) (storage.Event, error) {
	if err := a.s.Update(event); err != nil {
		return storage.Event{}, err
	}

	return event, nil
}

func (a *App) DeleteEvent(event storage.Event) error {
	if err := a.s.DeleteEvent(event); err != nil {
		return err
	}

	return nil
}

func (a *App) GetEvents(startAt time.Time, endAt time.Time, userID int) ([]storage.Event, error) {
	return a.s.GetEvents(startAt, endAt, userID)
}
