package mq

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/notificator"
)

type Consumer struct {
	c   sarama.Consumer
	n   Notificator
	l   Logger
	s   Storage
	cfg Config
}

type Notificator interface {
	SentNotification(notification notificator.Notification) error
}

type Storage interface {
	Notify(EventID string) error
}

func NewConsumer(c sarama.Consumer, n Notificator, l Logger, cfg Config, s Storage) *Consumer {
	return &Consumer{c: c, n: n, l: l, cfg: cfg, s: s}
}

func (c Consumer) Receive(ctx context.Context) error {
	defer func() {
		if err := c.c.Close(); err != nil {
			c.l.Error("Can`t close connect consumer: " + err.Error())
		}
	}()

	partitionConsumer, err := c.c.ConsumePartition(c.cfg.GetKafkaTopic(), 0, sarama.OffsetNewest)
	if err != nil {
		c.l.Error("Can`t create partition consumer: " + err.Error())
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			c.l.Error("Can`t close partition consumer: " + err.Error())
		}
	}()

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var n notificator.Notification
			if err := json.Unmarshal(msg.Value, &n); err != nil {
				c.l.Error("Can`t deserialize from: " + string(msg.Value) + ". Error: " + err.Error())

				continue
			}

			if err := c.n.SentNotification(n); err != nil {
				c.l.Error("Can`t sent notification: " + err.Error())

				continue
			}

			if err := c.s.Notify(n.EventID); err != nil {
				c.l.Error("Can`t update notification status in db: " + err.Error())

				continue
			}

			c.l.Info("Handle msg: " + string(msg.Value))
		case <-ctx.Done():
			c.l.Info("Stop consumer")
			break ConsumerLoop
		}
	}

	c.l.Info("Done consumer")
	return nil
}
