package main

import (
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/logger"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/rabbit"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/storage"
)

type Config struct {
	Logger  logger.Conf
	Storage storage.Conf
	Rabbit  rabbit.Config
}
