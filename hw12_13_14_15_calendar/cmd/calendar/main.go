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
	internalgrpc "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"gopkg.in/yaml.v2"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Init: App Config
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	logg, err := logger.New(config.Logger.File, config.Logger.Level, config.Logger.Formatter)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Init: Storage
	var storage app.Storage
	switch config.Storage.Type {
	case StorageMem:
		storage = memorystorage.New()
	case StorageSQL:
		sqlStorage := sqlstorage.New(ctx, config.Storage.Dsn, logg)
		sqlStorage.Connect(ctx)
		storage = sqlStorage
	default:
		log.Fatalf("Unknown storage type: %s\n", config.Storage.Type)
	}
	defer cancel()

	calendar := app.New(logg, storage)

	serverGrpc := internalgrpc.NewServer(config.GRPC.Port, calendar, logg)

	// Осторожно завершаем работу HTTP сервера
	go func() {
		<-ctx.Done()
		serverGrpc.Stop()
	}()

	go serverGrpc.Start()

	server := internalhttp.NewServer(logg, calendar, config.HTTP.Host, config.HTTP.Port)

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

func loadConfig(configPath string) (*Config, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %w", configPath, err)
	}

	cfg := NewConfig()
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read yaml: %w", err)
	}

	return &cfg, nil
}
