package main

import (
	"context"
	"flag"
	"log"
	"os"
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
	flag.StringVar(&configFile, "config", os.Getenv("CONFIG_FILE"), "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg := &Config{}
	if configFile != "" {
		err := configutil.LoadConfigFromFile(configFile, cfg)
		if err != nil {
			log.Fatalf("Failed to load config: %s", err)
		}
	} else {
		err := configutil.LoadConfigFromEnv(cfg)
		if err != nil {
			log.Fatalf("Failed to load config: %s", err)
		}
	}

	logg, err := logger.New(cfg.Logger.File, cfg.Logger.Level, cfg.Logger.Formatter)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	rcv, err := queue.NewRabbitQueue(ctx, cfg.Rabbit.Dsn, cfg.Rabbit.Exchange, cfg.Rabbit.Queue, logg)
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
