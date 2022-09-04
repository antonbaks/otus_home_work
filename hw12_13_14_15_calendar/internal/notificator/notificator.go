package notificator

import (
	"os"
	"time"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Notification struct {
	EventID      string
	EventTitle   string
	UserID       int
	EventStartAt time.Time
}

func (n Notification) String() string {
	return n.EventTitle
}

type Notificator struct {
	s Storage
}

type Storage interface {
	GetEventsForRemind(startAt time.Time, endAt time.Time) ([]storage.Event, error)
}

func NewNotificator(s Storage) *Notificator {
	return &Notificator{s: s}
}

func NewNotificatorWithoutDB() *Notificator {
	return &Notificator{}
}

func (n Notificator) GetNotifications(startAt time.Time, endAt time.Time) ([]Notification, error) {
	events, err := n.s.GetEventsForRemind(startAt, endAt)
	if err != nil {
		return nil, err
	}

	return createNotificationsByEvents(events), nil
}

func (n Notificator) SentNotification(notification Notification) error {
	if _, err := os.Stdout.WriteString("Send notification: " + notification.String() + "\n"); err != nil {
		return err
	}

	return nil
}

func createNotificationsByEvents(events []storage.Event) []Notification {
	notifications := make([]Notification, 0)
	for _, e := range events {
		notifications = append(notifications, Notification{
			EventID:      e.ID,
			EventTitle:   e.Title,
			EventStartAt: e.StartAt,
			UserID:       e.UserID,
		})
	}

	return notifications
}
