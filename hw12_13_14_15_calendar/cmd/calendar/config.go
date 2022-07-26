package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

const StorageInMemory = "in_memory"

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageTypeConf
	HTTP    HTTPConf
	SQL     SQLConf
}

type LoggerConf struct {
	Level string
}

type StorageTypeConf struct {
	Type string
}

type HTTPConf struct {
	Host string
	Port string
}

type SQLConf struct {
	DataSourceName string
	DriverName     string
	MigrationsDir  string
}

func NewConfig(path string) Config {
	var c Config

	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Config) GetHTTPHost() string {
	return c.HTTP.Host
}

func (c *Config) GetHTTPPort() string {
	return c.HTTP.Port
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
