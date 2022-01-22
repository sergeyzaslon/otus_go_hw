package main

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	HTTP    HTTPConf
	GRPC    GRPCConf
}

type LoggerConf struct {
	Level     string
	File      string
	Formatter string
}

const (
	StorageMem = "memory"
	StorageSQL = "sql"
)

type StorageConf struct {
	Type string
	Dsn  string
}

type HTTPConf struct {
	Host string
	Port int
}

type GRPCConf struct {
	Host string
	Port int
}

func NewConfig() Config {
	return Config{}
}
