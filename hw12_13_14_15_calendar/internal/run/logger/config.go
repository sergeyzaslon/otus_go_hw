package logger

type Conf struct {
	Level     string `env:"LOG_LEVEL" default:"info"`
	File      string `env:"LOG_FILE" default:"stderr"`
	Formatter string `env:"LOG_FORMAT" default:"json"`
}
