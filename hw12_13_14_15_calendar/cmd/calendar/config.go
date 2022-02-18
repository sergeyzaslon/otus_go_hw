package main

import (
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/logger"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/storage"
)

type Config struct {
	Logger  logger.Conf
	Storage storage.Conf
	HTTP    HTTPConf
	GRPC    GRPCConf
	File    string `env:"LOG_FILE" default:"stderr"`
}

type HTTPConf struct {
	Host string `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port int    `env:"HTTP_PORT" envDefault:"80"`
}

type GRPCConf struct {
	Host string `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port int    `env:"HTTP_PORT" envDefault:"8080"`
}
