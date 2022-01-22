package app

type Logger interface {
	Debug(msg string, params ...interface{})
	Info(msg string, params ...interface{})
	Warn(msg string, params ...interface{})
	Error(msg string, params ...interface{})
}
