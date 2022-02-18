package rabbit

type Config struct {
	Dsn      string `env:"QUEUE_DSN,required"`
	Queue    string `env:"RABBIT_QUEUE" envDefault:"0.0.0.0"`
	Exchange string `env:"RABBIT_EXCHANGE" envDefault:"0.0.0.0"`
}
