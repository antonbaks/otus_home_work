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
	internalgrpc "github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/internal/server/grpc"
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

	logg := logger.New(config.Logger.Level, os.Stderr, os.Stdout)

	var storage app.Storage
	if config.Storage.Type == StorageInMemory {
		storage = memorystorage.New(logg)
	} else {
		storage = sqlstorage.New(&config, logg)
		if err := storage.MigrationUp(); err != nil {
			log.Fatalln(err)
		}
	}

	calendar := app.New(logg, storage)

	service := internalgrpc.NewService(calendar)
	grpcServer := internalgrpc.NewServer(*service, logg, &config)
	gwServ := internalhttp.NewServer(logg, &config, grpcServer.ServMux)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		_, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := storage.Close(); err != nil {
			logg.Error("failed to stop db: " + err.Error())
		}

		if err := gwServ.Stop(); err != nil {
			logg.Error("failed to stop gw serv: " + err.Error())
		}

		grpcServer.Stop()
	}()

	logg.Info("calendar is running...")

	if err := grpcServer.Start(ctx); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		cancel()
		os.Exit(1) // nolint:gocritic
	}

	if err := gwServ.Start(); err != nil {
		logg.Error("failed to start gw server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
