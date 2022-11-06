package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/cleaner"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/mq"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/notificator"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/scheduler"
	sqlstorage "github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level, os.Stderr, os.Stdout)

	storage := sqlstorage.New(&config, logg)
	if err := storage.MigrationUp(); err != nil {
		logg.Error("failed connect db: " + err.Error())
		os.Exit(1)
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaProducer, err := sarama.NewSyncProducer(config.Kafka.Hosts, kafkaConfig)
	if err != nil {
		logg.Error("failed create kafka producer: " + err.Error())
		os.Exit(1)
	}

	n := notificator.NewNotificator(storage)
	p := mq.NewProducer(kafkaProducer, logg, &config)
	c := cleaner.NewCleaner(logg, storage, &config)

	newScheduler := scheduler.NewScheduler(&config, n, p, logg, c)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		if err := newScheduler.Stop(); err != nil {
			logg.Error("failed stop scheduler: " + err.Error())
		}

		if err := kafkaProducer.Close(); err != nil {
			logg.Error("failed stop kafka client: " + err.Error())
		}

		if err != nil {
			os.Exit(1)
		}
	}()

	logg.Info("calendar_scheduler is running...")
	if err := newScheduler.Start(); err != nil {
		logg.Error("failed start scheduler: " + err.Error())
		cancel()
		os.Exit(1) // nolint:gocritic
	}
}
