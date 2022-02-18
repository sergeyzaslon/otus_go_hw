package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

const (
	StorageMem = "memory"
	StorageSQL = "sql"
)

func Create(ctx context.Context, conf Conf, logger app.Logger) (app.Storage, error) {
	var storage app.Storage
	if strings.HasPrefix(conf.Dsn, "memory://") {
		storage = memorystorage.New()
	} else {
		sqlStorage := sqlstorage.New(ctx, conf.Dsn, logger)

		maxTries := 10
		tries := 0
		var err error
		for tries < 10 {
			err = sqlStorage.Connect(ctx)
			if err == nil {
				break
			}

			logger.Warn("Failed to connect to SQL storage: %w", err)
			time.Sleep(time.Second * 3)
			tries++
		}

		if tries == maxTries && err != nil {
			return nil, fmt.Errorf("failed to create SQL Storage: %w", err)
		}

		storage = sqlStorage
	}

	return storage, nil
}
