package scheduler

import (
	"fmt"
	"time"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/notificator"
)

type Scheduler struct {
	cfg  Config
	n    Notificator
	p    Producer
	l    Logger
	c    Cleaner
	done bool
}

type Config interface {
	GetSchedulerTimeout() string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Notificator interface {
	GetNotifications(startAt time.Time, endAt time.Time) ([]notificator.Notification, error)
}

type Producer interface {
	SendMessages(notifications []notificator.Notification) error
}

type Cleaner interface {
	ClearEvents() error
}

func NewScheduler(cfg Config, n Notificator, p Producer, l Logger, c Cleaner) *Scheduler {
	return &Scheduler{cfg: cfg, n: n, p: p, done: false, l: l, c: c}
}

func (sch *Scheduler) Start() error {
	timeout, err := time.ParseDuration(sch.cfg.GetSchedulerTimeout())
	if err != nil {
		return err
	}

	for !sch.done {
		now := time.Now()
		afterNow := now.Add(timeout)
		n, err := sch.n.GetNotifications(now, afterNow)
		if err != nil {
			sch.l.Error(err.Error())
			continue
		}

		sch.l.Info(fmt.Sprintf("Between %s and %s find %d notifications", now.String(), afterNow.String(), len(n)))

		if err := sch.p.SendMessages(n); err != nil {
			sch.l.Error(err.Error())
		}

		if err := sch.c.ClearEvents(); err != nil {
			sch.l.Error(err.Error())
		}

		time.Sleep(timeout)
	}

	sch.l.Info("done scheduler")

	return nil
}

func (sch *Scheduler) Stop() error {
	sch.done = true

	sch.l.Info("stop scheduler")

	return nil
}
