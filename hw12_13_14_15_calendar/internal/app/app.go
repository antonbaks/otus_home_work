package app

import (
	"context"
	"io"
	"net/http"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	log Logger
	s   Storage
}

type Logger interface { // TODO
}

type Storage interface {
	CreateEvent(e storage.Event) error
	DeleteEvent(e storage.Event) error
	Update(e storage.Event) error
	GetEventByID(id string) (storage.Event, error)
	GetAllEvents() ([]storage.Event, error)
	MigrationUp(ctx context.Context) error
	Close(ctx context.Context) error
}

func New(log Logger, s Storage) *App {
	return &App{
		log: log,
		s:   s,
	}
}

func (a *App) CreateEvent(id, title string) error {
	return a.s.CreateEvent(storage.Event{ID: id, Title: title})
}

func (a *App) PrintHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, "Hello world")
}
