package mq

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/notificator"
)

type Producer struct {
	p sarama.SyncProducer
	c Config
	l Logger
}

func NewProducer(p sarama.SyncProducer, l Logger, c Config) *Producer {
	return &Producer{p: p, l: l, c: c}
}

type Config interface {
	GetKafkaTopic() string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func (p Producer) Send(n notificator.Notification) error {
	defer func() {
		if err := p.p.Close(); err != nil {
			p.l.Error("Can`t close connect producer: " + err.Error())
		}
	}()

	str, err := json.Marshal(n)
	if err != nil {
		return err
	}

	_, _, err = p.p.SendMessage(&sarama.ProducerMessage{
		Topic: p.c.GetKafkaTopic(),
		Value: sarama.StringEncoder(str),
	})

	if err != nil {
		p.l.Error("failed to send message to topic" + err.Error())
	}

	return nil
}
