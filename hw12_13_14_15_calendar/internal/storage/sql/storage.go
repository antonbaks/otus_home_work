package sqlstorage

import (
	"context"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	// Register db driver.
	_ "github.com/lib/pq"
	goose "github.com/pressly/goose/v3"
)

type Storage struct {
	cfg Config
	log Logger
	db  *sqlx.DB
}

type Config interface {
	GetDriverName() string
	GetDataSourceName() string
	GetMigrationDir() string
}

type Logger interface {
	Error(msg string)
}

func New(cfg Config, log Logger) *Storage {
	return &Storage{
		cfg: cfg,
		log: log,
	}
}

func (s *Storage) MigrationUp(ctx context.Context) error {
	if err := s.Connect(ctx); err != nil {
		return err
	}

	if err := goose.SetDialect(s.cfg.GetDriverName()); err != nil {
		return err
	}

	if err := goose.Up(s.db.DB, s.cfg.GetMigrationDir()); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Connect(s.cfg.GetDriverName(), s.cfg.GetDataSourceName())
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateEvent(e storage.Event) error {
	query := `
			INSERT INTO
			    events (title, description, start_at, end_at, user_id, remind_for)
			VALUES
			    (:name, :description, :start_at, :end_at, :user_id, :remind_for)
	`

	if _, err := s.db.NamedExec(query, e); err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}

func (s *Storage) DeleteEvent(e storage.Event) error {
	query := `
			DELETE FROM 
				events 
		   	WHERE 
		   	    id=:id
	`

	if _, err := s.db.NamedExec(query, e); err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}

func (s *Storage) Update(e storage.Event) error {
	query := `
			UPDATE 
			    events 
			SET
			    title=:title,
			    description=:description,
			    start_at=:start_at,
			    end_at=:end_at,
			    user_id=:user_id,
			    remind_for=:remind_for
			WHERE
			    id=:id
	`

	if _, err := s.db.NamedExec(query, e); err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}

func (s *Storage) GetEventByID(id string) (storage.Event, error) {
	query := `
			SELECT 
				title,
				description,
				start_at,
				end_at,
				user_id,
				remind_for
			FROM
				events
			WHERE
				id=$1
	`

	e := storage.Event{}
	if err := s.db.Get(&e, query, id); err != nil {
		return e, err
	}

	return e, nil
}

func (s *Storage) GetAllEvents() ([]storage.Event, error) {
	query := `
			SELECT 
				title,
				description,
				start_at,
				end_at,
				user_id,
				remind_for
			FROM
				events
	`

	var e []storage.Event
	if err := s.db.Select(&e, query); err != nil {
		return e, err
	}

	return e, nil
}
