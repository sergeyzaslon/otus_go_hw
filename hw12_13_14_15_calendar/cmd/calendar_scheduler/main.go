package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/clock"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/configutil"
	storagefactory "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/storage"
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

	storage, err := storagefactory.Create(ctx, config.Storage, logg)
	if err != nil {
		log.Fatalf("Failed to create storage: %s", err)
	}

	var (
		eventSource app.EventSource
		ok          bool
	)
	if eventSource, ok = storage.(app.EventSource); !ok {
		log.Fatalf("Storage does not implement app.EventSource interface")
	}

	rcv, err := queue.NewRabbitQueue(ctx, config.Rabbit.Dsn, config.Rabbit.Exchange, config.Rabbit.Queue, logg)
	if err != nil {
		log.Fatalf("Failed to create NotificationSender (rabbit): %s", err)
	}

	clck := clock.NewSystemClock()

	scheduler := app.NewAppScheduler(eventSource, rcv, clck, logg)

	timer := time.Tick(time.Second)
	timerHour := time.Tick(time.Hour)

	go func() {
		for {
			select {
			case <-timer:
				err := scheduler.Notify()
				if err != nil {
					logg.Error("Failed to Notify: %s", err)
				}
			case <-timerHour:
				err := scheduler.RemoveOneYearOld()
				if err != nil {
					logg.Error("Failed to Notify: %s", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	logg.Info("Calendar Scheduler Started!")

	<-ctx.Done()
}
