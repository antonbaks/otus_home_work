package cleaner

import (
	"time"
)

type Clean struct {
	TimeAfter time.Time `db:"end_at"`
}

type Cleaner struct {
	l Logger
	s Storage
	c Config
}

type Logger interface {
	Info(msg string)
}

type Config interface {
	GetCleanAfter() string
}

type Storage interface {
	DeleteByEndAt(clean Clean) error
}

func NewCleaner(l Logger, s Storage, c Config) *Cleaner {
	return &Cleaner{l: l, s: s, c: c}
}

func (c Cleaner) ClearEvents() error {
	cleanAfter, err := time.ParseDuration(c.c.GetCleanAfter())
	if err != nil {
		return err
	}

	timeAfter := time.Now().Add(cleanAfter)
	if err := c.s.DeleteByEndAt(Clean{TimeAfter: timeAfter}); err != nil {
		return err
	}

	c.l.Info("Clean events after: " + timeAfter.String())

	return nil
}
