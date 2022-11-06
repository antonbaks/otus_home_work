package sender

import "context"

type Sender struct {
	c Consumer
	l Logger
}

type Consumer interface {
	Receive(ctx context.Context) error
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func NewSender(l Logger, c Consumer) *Sender {
	return &Sender{c: c, l: l}
}

func (s Sender) Start(ctx context.Context) error {
	if err := s.c.Receive(ctx); err != nil {
		s.l.Error("Sender start with error: " + err.Error())
		return err
	}

	return nil
}
