package storage

type Conf struct {
	Dsn string `env:"STORAGE_DSN,required"`
}
