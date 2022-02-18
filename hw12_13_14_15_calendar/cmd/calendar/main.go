package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/configutil"
	storagefactory "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/run/storage"
	internalgrpc "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", os.Getenv("CONFIG_FILE"), "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Init: App Config
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

	fmt.Println(cfg)

	logg, err := logger.New(cfg.Logger.File, cfg.Logger.Level, cfg.Logger.Formatter)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	storage, err := storagefactory.Create(ctx, cfg.Storage, logg)
	if err != nil {
		log.Fatalf("Failed to create storage: %s", err)
	}

	calendar := app.New(logg, storage)

	serverGrpc := internalgrpc.NewServer(cfg.GRPC.Port, calendar, logg)

	// Осторожно завершаем работу HTTP сервера
	go func() {
		<-ctx.Done()
		serverGrpc.Stop()
	}()

	go serverGrpc.Start()

	server := internalhttp.NewServer(logg, calendar, cfg.HTTP.Host, cfg.HTTP.Port)

	// Осторожно завершаем работу HTTP сервера
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("Failed to stop http server: " + err.Error())
		}
	}()

	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("Failed to start http server: " + err.Error())
		}
	}()

	<-ctx.Done()
}
