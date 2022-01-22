package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/configutil"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/transport"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	// Init: App Config
	config := &Config{}
	err := configutil.LoadConfig(configFile, config)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	logg, err := logger.New(config.Logger.File, config.Logger.Level, config.Logger.Formatter)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	rcv, err := queue.NewRabbitQueue(ctx, config.Rabbit.Dsn, config.Rabbit.Exchange, config.Rabbit.Queue, logg)
	if err != nil {
		log.Fatalf("Filed to create NotificationSender (rabbit): %s", err)
	}

	transports := []app.NotificationTransport{
		transport.NewLogNotificationTransport(logg),
	}

	scheduler := app.NewNotificationSender(rcv, logg, transports)

	scheduler.Run()

	<-ctx.Done()
}
