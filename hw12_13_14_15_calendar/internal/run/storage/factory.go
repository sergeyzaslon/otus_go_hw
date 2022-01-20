package storage

import (
	"context"
	"fmt"

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
	switch conf.Type {
	case StorageMem:
		storage = memorystorage.New()
	case StorageSQL:
		sqlStorage := sqlstorage.New(ctx, conf.Dsn, logger)
		sqlStorage.Connect(ctx)
		storage = sqlStorage
	default:
		return nil, fmt.Errorf("unknown storage type %s, available: [%s, %s]", conf.Type, StorageMem, StorageSQL)
	}

	return storage, nil
}
