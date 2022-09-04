package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/mq"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/notificator"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/sender"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level, os.Stderr, os.Stdout)

	kafkaConfig := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(config.Kafka.Hosts, kafkaConfig)
	if err != nil {
		logg.Error(err.Error())
		os.Exit(1)
	}

	n := notificator.NewNotificatorWithoutDB()
	mqConsumer := mq.NewConsumer(consumer, n, logg, &config)
	newSender := sender.NewSender(logg, mqConsumer)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar_sender is running...")

	if err := newSender.Start(ctx); err != nil {
		logg.Error("calendar_sender start with error " + err.Error())
		cancel()
		os.Exit(1) // nolint:gocritic
	}
}
