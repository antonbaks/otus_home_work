package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger    LoggerConf
	SQL       SQLConf
	Scheduler SchedulerConf
	Cleaner   CleanerConf
	Kafka     KafkaConf
}

type LoggerConf struct {
	Level string
}

type SQLConf struct {
	DataSourceName string
	DriverName     string
	MigrationsDir  string
}

type SchedulerConf struct {
	Timeout string
}

type CleanerConf struct {
	CleanAfter string
}

type KafkaConf struct {
	Topic string
	Hosts []string
}

func NewConfig(path string) Config {
	var c Config

	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Config) GetDriverName() string {
	return c.SQL.DriverName
}

func (c *Config) GetDataSourceName() string {
	return c.SQL.DataSourceName
}

func (c *Config) GetMigrationDir() string {
	return c.SQL.MigrationsDir
}

func (c *Config) GetSchedulerTimeout() string {
	return c.Scheduler.Timeout
}

func (c *Config) GetCleanAfter() string {
	return c.Cleaner.CleanAfter
}

func (c *Config) GetKafkaTopic() string {
	return c.Kafka.Topic
}
