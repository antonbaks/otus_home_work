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

func (p Producer) SendMessages(n []notificator.Notification) error {
	producerMessage := make([]*sarama.ProducerMessage, len(n))

	for i, oneN := range n {
		str, err := json.Marshal(oneN)
		if err != nil {
			return err
		}

		producerMessage[i] = &sarama.ProducerMessage{
			Topic: p.c.GetKafkaTopic(),
			Value: sarama.StringEncoder(str),
		}
	}

	if err := p.p.SendMessages(producerMessage); err != nil {
		p.l.Error("failed to send messages to topic" + err.Error())

		return err
	}

	return nil
}
