package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)

	logg := logger.New(config.Logger.Level)

	var storage app.Storage
	if config.Storage.Type == StorageInMemory {
		storage = memorystorage.New(logg)
	} else {
		storage = sqlstorage.New(&config, logg)
		if err := storage.MigrationUp(context.Background()); err != nil {
			log.Fatalln(err)
		}
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, &config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err := storage.Close(ctx); err != nil {
			logg.Error("failed to stop db: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
