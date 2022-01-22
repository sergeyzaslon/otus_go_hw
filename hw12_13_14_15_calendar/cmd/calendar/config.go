package main

import (
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/logger"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/storage"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  logger.Conf
	Storage storage.Conf
	HTTP    HTTPConf
	GRPC    GRPCConf
}

type HTTPConf struct {
	Host string
	Port int
}

type GRPCConf struct {
	Host string
	Port int
}
